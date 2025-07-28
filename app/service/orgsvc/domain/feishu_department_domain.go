package orgsvc

import (
	"fmt"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdkConsts "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/errors"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func DeptAdd(orgId int64, deptOpenId string, sourceChannel, corpId string) (int64, errs.SystemErrorInfo) {
	deptInfo := &po.PpmOrgDepartmentOutInfo{}
	err := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcIsDelete:           consts.AppIsNoDelete,
		consts.TcOrgId:              orgId,
		consts.TcOutOrgDepartmentId: deptOpenId,
	}, deptInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			//没有就新增
			deptId, addDeptErr := addDept(orgId, deptOpenId, sourceChannel, corpId)
			if addDeptErr != nil {
				log.Error(addDeptErr)
				return 0, addDeptErr
			}

			return deptId, nil
		} else {
			log.Error(err)
			return 0, errs.MysqlOperateError
		}
	}

	//如果已经有了就直接返回
	return deptInfo.DepartmentId, nil
}

func addDept(orgId int64, deptOpenId string, sourceChannel, corpId string) (int64, errs.SystemErrorInfo) {
	//先上锁
	lockKey := consts.AddFeishuDepartmentLock + fmt.Sprintf("%d:%s", orgId, deptOpenId)
	uuid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uuid)
	if lockErr != nil {
		log.Error(lockErr)
		return 0, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if !suc {
		log.Errorf("新建部门时没有获取到锁 orgId %d departmentId %v", orgId, deptOpenId)
		return 0, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	defer func() {
		if _, err := cache.ReleaseDistributedLock(lockKey, uuid); err != nil {
			log.Error(err)
		}
	}()

	deptInfo := &po.PpmOrgDepartmentOutInfo{}
	deptInfoErr := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcIsDelete:           consts.AppIsNoDelete,
		consts.TcOrgId:              orgId,
		consts.TcOutOrgDepartmentId: deptOpenId,
	}, deptInfo)
	if deptInfoErr == nil {
		//如果已经创建就直接返回
		return deptInfo.Id, nil
	} else {
		if deptInfoErr != db.ErrNoMoreRows {
			log.Error(deptInfoErr)
			return 0, errs.MysqlOperateError
		}
	}

	client, err := platform_sdk.GetClient(sourceChannel, corpId)
	if err != nil {
		log.Errorf("[GetClient] err:%v", err)
		return 0, errs.PlatFormOpenApiCallError
	}
	departmentInfo, sdkErr := client.GetDepartmentInfo(&sdkVo.GetDepartmentInfoReq{DeptId: deptOpenId})
	if sdkErr != nil {
		// 如果是无效部门，跳过处理
		if errors.CheckDepartmentInvalidError(sdkErr) {
			return 0, nil
		}
		log.Errorf("[GetDepartmentInfo] deptId:%v, err:%v", deptOpenId, sdkErr)
		return 0, errs.PlatFormOpenApiCallError
	}

	if departmentInfo.Info.Status == sdkConsts.DepartmentStatusInvalid {
		//部门状态，0 无效，1 有效
		log.Infof("部门删除后提供的一次部门更新回调:%s", deptOpenId)
		return 0, nil
	}

	outId, outIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableDepartmentOutInfo)
	if outIdErr != nil {
		log.Error(outIdErr)
		return 0, outIdErr
	}
	deptId, deptIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableDepartment)
	if deptIdErr != nil {
		log.Error(deptIdErr)
		return 0, deptIdErr
	}

	parentId := int64(0)
	if departmentInfo.Info.ParentId != "" && departmentInfo.Info.ParentId != "0" {
		//父任务
		parentInfo := &po.PpmOrgDepartmentOutInfo{}
		parentInfoErr := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
			consts.TcIsDelete:           consts.AppIsNoDelete,
			consts.TcOrgId:              orgId,
			consts.TcOutOrgDepartmentId: departmentInfo.Info.ParentId,
		}, parentInfo)
		if parentInfoErr != nil {
			if parentInfoErr == db.ErrNoMoreRows {
				//父任务没有导入到系统 就去新增父任务
				initParentId, err := addDept(orgId, departmentInfo.Info.ParentId, sourceChannel, corpId)
				if err != nil {
					log.Error(err)
					return 0, err
				}
				parentId = initParentId
			} else {
				log.Error(parentInfoErr)
				return 0, errs.MysqlOperateError
			}
		} else {
			parentId = parentInfo.DepartmentId
		}
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := mysql.TransInsert(tx, &po.PpmOrgDepartment{
			Id:             deptId,
			OrgId:          orgId,
			Name:           departmentInfo.Info.Name,
			ParentId:       parentId,
			SourcePlatform: "",
			SourceChannel:  sourceChannel,
		})
		if err != nil {
			log.Error(err)
			return err
		}
		err1 := mysql.TransInsert(tx, &po.PpmOrgDepartmentOutInfo{
			Id:                       outId,
			OrgId:                    orgId,
			DepartmentId:             deptId,
			SourcePlatform:           "",
			SourceChannel:            sourceChannel,
			OutOrgDepartmentId:       departmentInfo.Info.Id,
			OutOrgDepartmentCode:     "",
			Name:                     departmentInfo.Info.Name,
			OutOrgDepartmentParentId: departmentInfo.Info.ParentId,
		})
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return 0, errs.MysqlOperateError
	}

	return deptId, nil
}

func DeptDelete(orgId int64, deptOpenId string) (int64, errs.SystemErrorInfo) {
	deptInfo := &po.PpmOrgDepartmentOutInfo{}
	err := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcIsDelete:           consts.AppIsNoDelete,
		consts.TcOrgId:              orgId,
		consts.TcOutOrgDepartmentId: deptOpenId,
	}, deptInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			//没有就正常返回
			return 0, nil
		} else {
			log.Error(err)
			return 0, errs.MysqlOperateError
		}
	}

	//如果已经有了就删除
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//用户部门关联
		_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcDepartmentId: deptInfo.DepartmentId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if err1 != nil {
			log.Error(err1)
			return err1
		}
		//外部信息
		_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{
			consts.TcId: deptInfo.Id,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if err2 != nil {
			log.Error(err2)
			return err2
		}
		//主表
		_, err3 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
			consts.TcId: deptInfo.DepartmentId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if err3 != nil {
			log.Error(err3)
			return err3
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return 0, errs.MysqlOperateError
	}

	return deptInfo.DepartmentId, nil
}

// 更新部门（更新后的部门的子部门也许不存在 需要新增 todo）
func DeptUpdate(orgId int64, deptOpenId string, sourceChannel, corpId string) (int64, errs.SystemErrorInfo) {
	client, sdkErr := platform_sdk.GetClient(sourceChannel, corpId)
	if sdkErr != nil {
		log.Errorf("[GetClient] err:%v", sdkErr)
		return 0, errs.PlatFormOpenApiCallError
	}

	departmentInfoReply, sdkErr := client.GetDepartmentInfo(&sdkVo.GetDepartmentInfoReq{DeptId: deptOpenId})
	if sdkErr != nil {
		log.Errorf("[GetDepartmentInfo] deptId:%v, err:%v", deptOpenId, sdkErr)
		return 0, errs.PlatFormOpenApiCallError
	}
	departmentInfo := departmentInfoReply.Info
	if departmentInfo.Status == sdkConsts.DepartmentStatusInvalid {
		//部门状态，0 无效，1 有效
		log.Infof("部门删除后提供的一次部门更新回调:%s", deptOpenId)
		return 0, nil
	}

	deptInfo := &po.PpmOrgDepartmentOutInfo{}
	err := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcIsDelete:           consts.AppIsNoDelete,
		consts.TcOrgId:              orgId,
		consts.TcOutOrgDepartmentId: deptOpenId,
	}, deptInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			//没有就创建（可能之前没有同步）
			deptId, addDeptErr := addDept(orgId, deptOpenId, sourceChannel, corpId)
			if addDeptErr != nil {
				log.Error(addDeptErr)
				return 0, addDeptErr
			}
			return deptId, nil
		} else {
			log.Error(err)
			return 0, errs.MysqlOperateError
		}
	}

	//查看名称和父部门是否变化，如果没有变化直接返回就好了
	if departmentInfo.Name == deptInfo.Name && departmentInfo.ParentId == deptInfo.OutOrgDepartmentParentId {
		return deptInfo.DepartmentId, nil
	}

	upd := mysql.Upd{}
	truUpd := mysql.Upd{}
	if departmentInfo.Name != deptInfo.Name {
		upd[consts.TcName] = departmentInfo.Name
		truUpd[consts.TcName] = departmentInfo.Name
	}
	if departmentInfo.ParentId != deptInfo.OutOrgDepartmentParentId {
		upd[consts.TcOutOrgDepartmentParentId] = departmentInfo.ParentId
		parentId := int64(0)
		//父部门
		if departmentInfo.ParentId != "" {
			parentInfo := &po.PpmOrgDepartmentOutInfo{}
			parentInfoErr := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
				consts.TcIsDelete:           consts.AppIsNoDelete,
				consts.TcOrgId:              orgId,
				consts.TcOutOrgDepartmentId: departmentInfo.ParentId,
			}, parentInfo)
			if parentInfoErr != nil {
				if parentInfoErr == db.ErrNoMoreRows {
					//父部门没有导入到系统 就去新增父部门
					initParentId, err := addDept(orgId, departmentInfo.ParentId, sourceChannel, corpId)
					if err != nil {
						log.Error(err)
						return 0, err
					}
					parentId = initParentId
				} else {
					log.Error(parentInfoErr)
					return 0, errs.MysqlOperateError
				}
			} else {
				parentId = parentInfo.DepartmentId
			}
		}

		truUpd[consts.TcParentId] = parentId
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{
			consts.TcId: deptInfo.Id,
		}, upd)
		if err != nil {
			log.Error(err)
			return err
		}

		_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
			consts.TcId: deptInfo.DepartmentId,
		}, truUpd)
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return 0, errs.MysqlOperateError
	}

	return deptInfo.DepartmentId, nil
}

// 部门变更，包括最新的部门层级
func ChangeDeptScope(orgId int64, sourceChannel, corpId string) errs.SystemErrorInfo {
	//查询原有的所有部门
	originDepts := &[]po.PpmOrgDepartmentOutInfo{}
	originDeptsErr := mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, originDepts)
	if originDeptsErr != nil {
		log.Error(originDeptsErr)
		return errs.MysqlOperateError
	}

	//在不在授权范围通过status来决定
	notScopeDeptIds := []string{}
	scopeDeptIds := []string{}
	outDeptToIdMap := map[string]int64{} //外部id对应内部id

	for _, info := range *originDepts {
		if info.OutOrgDepartmentId != "0" {
			outDeptToIdMap[info.OutOrgDepartmentId] = info.DepartmentId
		}
		if info.Status == 1 {
			scopeDeptIds = append(scopeDeptIds, info.OutOrgDepartmentId)
		} else {
			notScopeDeptIds = append(notScopeDeptIds, info.OutOrgDepartmentId)
		}
	}

	client, sdkErr := platform_sdk.GetClient(sourceChannel, corpId)
	if sdkErr != nil {
		log.Errorf("[GetClient] err:%v", sdkErr)
		return errs.PlatFormOpenApiCallError
	}

	//获取授权范围
	deptsReply, sdkErr := client.GetScopeDeps()
	if sdkErr != nil {
		log.Errorf("[GetScopeDeps] err:%v", sdkErr)
		return errs.PlatFormOpenApiCallError
	}
	deptMap := make(map[string]*sdkVo.DepartmentInfo, len(deptsReply.Depts))
	trulyScopeIds := make([]string, 0, len(deptsReply.Depts))
	for _, info := range deptsReply.Depts {
		trulyScopeIds = append(trulyScopeIds, info.Id)
		deptMap[info.Id] = info
	}

	//系统授权范围和真正授权范围比较（需要删除|初步需要增加的）
	del, midAdd := util.GetDifMemberIdsByString(scopeDeptIds, trulyScopeIds)
	upd := []string{}
	add := []string{}
	for _, s := range midAdd {
		if ok, _ := slice.Contain(notScopeDeptIds, s); ok {
			//更新为授权可用
			upd = append(upd, s)
		} else {
			//新增
			add = append(add, s)
		}
	}

	delDeptIds := []int64{}
	updDeptIds := []int64{}
	for _, info := range *originDepts {
		if ok, _ := slice.Contain(del, info.OutOrgDepartmentId); ok {
			delDeptIds = append(delDeptIds, info.DepartmentId)
		}
		if ok, _ := slice.Contain(upd, info.OutOrgDepartmentId); ok {
			updDeptIds = append(updDeptIds, info.DepartmentId)
		}
	}

	outInsertPos := []po.PpmOrgDepartmentOutInfo{}
	insertPos := []po.PpmOrgDepartment{}
	if len(add) > 0 {
		ids, idsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartment, len(add))
		if idsErr != nil {
			log.Error(idsErr)
			return idsErr
		}
		outIds, outIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartmentOutInfo, len(add))
		if outIdsErr != nil {
			log.Error(outIdsErr)
			return outIdsErr
		}
		for i, s := range add {
			//把接下来要新增的部门 也转化成 map， 这样 就把所有相关部门都取出来了，从而获取所有父部门map
			//if deptMap[s].ParentId != "" {
			//	if _, ok := outDeptToIdMap[deptMap[s].ParentId]; !ok {
			//		outDeptToIdMap[deptMap[s].Id] = ids.Ids[i].Id
			//	}
			//}
			if _, ok := outDeptToIdMap[s]; !ok {
				outDeptToIdMap[deptMap[s].Id] = ids.Ids[i].Id
			}
		}

		for i, s := range add {
			parentId := int64(0)
			if id, ok := outDeptToIdMap[deptMap[s].ParentId]; ok {
				parentId = id
			}
			insertPos = append(insertPos, po.PpmOrgDepartment{
				Id:             ids.Ids[i].Id,
				OrgId:          orgId,
				Name:           deptMap[s].Name,
				ParentId:       parentId,
				SourcePlatform: sourceChannel,
				SourceChannel:  "",
				Status:         1,
			})
			outInsertPos = append(outInsertPos, po.PpmOrgDepartmentOutInfo{
				Id:                       outIds.Ids[i].Id,
				OrgId:                    orgId,
				DepartmentId:             ids.Ids[i].Id,
				SourcePlatform:           sourceChannel,
				SourceChannel:            "",
				OutOrgDepartmentId:       deptMap[s].Id,
				OutOrgDepartmentCode:     "",
				Name:                     deptMap[s].Name,
				OutOrgDepartmentParentId: deptMap[s].ParentId,
				Status:                   1,
			})
		}
	}

	//处理要更新的 和 原来就有的 部门 （确认新的部门层级，可能部门层级发生了变动由于没有回调导致 部门层级不对）
	needUpdateParentMap := map[int64]string{} //需要更新父部门的  部门id->外部父部门id
	for _, info := range *originDepts {
		if ok, _ := slice.Contain(del, info.OutOrgDepartmentId); ok {
			//移除的不用判断
			continue
		}
		if outDeptInfo, ok := deptMap[info.OutOrgDepartmentId]; ok {
			if outDeptInfo.ParentId != info.OutOrgDepartmentParentId {
				//如果飞书查出来的父部门id != 本地的父部门id
				needUpdateParentMap[info.DepartmentId] = outDeptInfo.ParentId
			}
		}
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(add) > 0 {
			err1 := mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, slice.ToSlice(insertPos))
			if err1 != nil {
				log.Error(err1)
				return err1
			}

			err2 := mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, slice.ToSlice(outInsertPos))
			if err2 != nil {
				log.Error(err2)
				return err2
			}
		}

		if len(del) > 0 {
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{
				consts.TcOrgId:              orgId,
				consts.TcOutOrgDepartmentId: db.In(del),
			}, mysql.Upd{
				consts.TcStatus: 2,
			})
			if err1 != nil {
				log.Error(err1)
				return err1
			}

			_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcId:    db.In(delDeptIds),
			}, mysql.Upd{
				consts.TcStatus: 2,
			})
			if err2 != nil {
				log.Error(err2)
				return err2
			}
		}

		if len(upd) > 0 {
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{
				consts.TcOrgId:              orgId,
				consts.TcOutOrgDepartmentId: db.In(upd),
			}, mysql.Upd{
				consts.TcStatus: 1,
			})
			if err1 != nil {
				log.Error(err1)
				return err1
			}

			_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcId:    db.In(updDeptIds),
			}, mysql.Upd{
				consts.TcStatus: 1,
			})
			if err2 != nil {
				log.Error(err2)
				return err2
			}
		}
		if len(needUpdateParentMap) > 0 {
			for i, s := range needUpdateParentMap {
				_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{
					consts.TcOrgId:        orgId,
					consts.TcDepartmentId: i,
					consts.TcIsDelete:     consts.AppIsNoDelete,
				}, mysql.Upd{
					consts.TcOutOrgDepartmentParentId: s,
				})
				if err1 != nil {
					log.Error(err1)
					return err1
				}

				parentId := int64(0)
				if id, ok := outDeptToIdMap[s]; ok {
					parentId = id
				}
				_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
					consts.TcOrgId:    orgId,
					consts.TcId:       i,
					consts.TcIsDelete: consts.AppIsNoDelete,
				}, mysql.Upd{
					consts.TcParentId: parentId,
				})
				if err2 != nil {
					log.Error(err2)
					return err2
				}
			}
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.MysqlOperateError
	}
	log.Infof("更新部门情况:新增【%d】，删除【%d】，更新【%d】", len(add), len(del), len(upd))
	// 在同步用户和部门的关联关系前，检查是否有多余的外部部门信息，如果有，则删掉多余的。
	if err := CheckAndClearExtraOutDept(orgId); err != nil {
		log.Errorf("[ChangeDeptScope] orgId: %d, err: %v", orgId, err)
		return err
	}
	if err := ClearDeptForSyncDept(orgId); err != nil {
		log.Errorf("[ChangeDeptScope] orgId: %d, ClearDeptForSyncDept err: %v", orgId, err)
		return err
	}
	// 同步用户和部门的关联关系
	updUserDeptErr := handleSyncUserDept(client, orgId, corpId)
	if updUserDeptErr != nil {
		log.Error(updUserDeptErr)
		return updUserDeptErr
	}

	return nil
}

// CheckAndClearExtraOutDept 检查是否有多余的外部部门信息，如果有，则删掉多余的。 todo
func CheckAndClearExtraOutDept(orgId int64) errs.SystemErrorInfo {
	// 查询同一个 out deptId 有多条的情况
	outDeptStatArr := []po.DepartmentOutInfoGroupByOutDeptId{}
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf("[checkAndClearExtraOutDept] orgId: %d, err: %v", orgId, err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	err = conn.Collection(consts.TableDepartmentOutInfo).Find(cond).Select(db.Raw("out_org_department_id, name, count(*) cnt")).Group("out_org_department_id").All(&outDeptStatArr)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	if len(outDeptStatArr) < 1 {
		return nil
	}
	// 有重复数据的 outDeptId 列表
	hasDuplicatedOutDeptIdArr := make([]string, 0)
	for _, item := range outDeptStatArr {
		if item.Cnt <= 1 {
			continue
		}
		hasDuplicatedOutDeptIdArr = append(hasDuplicatedOutDeptIdArr, item.OutOrgDepartmentId)
	}
	if len(hasDuplicatedOutDeptIdArr) < 1 {
		return nil
	}
	// 查询出这批有多余的数据
	outDeptArr := []po.PpmOrgDepartmentOutInfo{}
	cond = db.Cond{
		consts.TcOrgId:              orgId,
		consts.TcIsDelete:           consts.AppIsNoDelete,
		consts.TcOutOrgDepartmentId: db.In(hasDuplicatedOutDeptIdArr),
	}
	err = conn.Collection(consts.TableDepartmentOutInfo).Find(cond).OrderBy("id desc").All(&outDeptArr)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	// 找出对于的数据的 id，并根据 id 删除
	delIds := make([]int64, 0)
	duplicateMap1 := make(map[string]struct{}, len(hasDuplicatedOutDeptIdArr))
	for _, outDept := range outDeptArr {
		if _, ok := duplicateMap1[outDept.OutOrgDepartmentId]; !ok {
			duplicateMap1[outDept.OutOrgDepartmentId] = struct{}{}
		} else {
			delIds = append(delIds, outDept.Id)
		}
	}
	if len(delIds) < 1 {
		return nil
	}
	// 执行删除
	// conn.SetLogging(true)
	delCond := db.Cond{
		consts.TcOrgId: orgId,
		consts.TcId:    db.In(delIds),
	}
	if err := conn.Collection(consts.TableDepartmentOutInfo).Find(delCond).Update(map[string]interface{}{
		consts.TcIsDelete: consts.AppIsDeleted,
	}); err != nil {
		log.Errorf("[checkAndClearExtraOutDept] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	log.Infof("[CheckAndClearExtraOutDept] 去除多余 out dept，ids: %s", json.ToJsonIgnoreError(delIds))

	return nil
}

// ClearDeptForSyncDept 同步飞书部门时，清理掉多余的内部部门数据 todo
func ClearDeptForSyncDept(orgId int64) errs.SystemErrorInfo {
	delCond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.NotIn(db.Raw("select department_id from "+consts.TableDepartmentOutInfo+" where is_delete=2 and org_id=?", orgId)),
	}
	cnt, err := mysql.UpdateSmartWithCond(consts.TableDepartment, delCond, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcVersion:  2022042801,
	})
	if err != nil {
		log.Errorf("[ClearDeptForSyncDept] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	log.Infof("[ClearDeptForSyncDept] 清理了 department 表，orgId: %d, 数据条数：[%d]", orgId, cnt)

	return nil
}

// 同步用户部门(有可能会丢失部分用户部门关联关系，因为人员同步是单个放到队列里面的，可能人员还没有更新完毕，就开始同步关联关系，不过影响甚小，可以让用户再次更新通讯录范围。)
func handleSyncUserDept(client sdk_interface.Sdk, orgId int64, tenantKey string) errs.SystemErrorInfo {
	//取出所有的用户部门关联
	userDeptInfo := &[]po.PpmOrgUserDepartment{}
	err := mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userDeptInfo)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	userDeptMap := map[int64][]int64{}
	// 一个用户在一个部门中的关联关系可能因为脏数据导致存在多条，为了能清理多余的所有数据，关联 id 为切片
	userDeptRelationMap := map[int64]map[int64][]int64{} //用户->部门->关联id
	for _, department := range *userDeptInfo {
		userDeptMap[department.UserId] = append(userDeptMap[department.UserId], department.DepartmentId)
		if _, ok := userDeptRelationMap[department.UserId]; ok {
			if len(userDeptRelationMap[department.UserId][department.DepartmentId]) < 1 {
				userDeptRelationMap[department.UserId][department.DepartmentId] = []int64{department.Id}
			} else {
				userDeptRelationMap[department.UserId][department.DepartmentId] =
					append(userDeptRelationMap[department.UserId][department.DepartmentId], department.Id)
			}
		} else {
			tempMap := map[int64][]int64{department.DepartmentId: []int64{department.Id}}
			userDeptRelationMap[department.UserId] = tempMap
		}
	}

	//获取所有外部部门信息
	outDeptInfo := &[]po.PpmOrgDepartmentOutInfo{}
	outDeptInfoErr := mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	}, outDeptInfo)
	if outDeptInfoErr != nil {
		log.Error(outDeptInfoErr)
		return errs.MysqlOperateError
	}
	outToInnerDeptId := map[string]int64{}
	if len(*outDeptInfo) == 0 {
		return nil
	}
	outDeptInfos := make([]*sdkVo.DepartmentInfo, 0, len(*outDeptInfo))
	for _, info := range *outDeptInfo {
		outToInnerDeptId[info.OutOrgDepartmentId] = info.DepartmentId
		outDeptInfos = append(outDeptInfos, &sdkVo.DepartmentInfo{Id: info.OutOrgDepartmentId})
	}

	//获取所有组织的用户外部信息
	outUserInfo := &[]po.PpmOrgUserOutInfo{}
	err1 := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	}, outUserInfo)
	if err1 != nil {
		log.Error(err1)
		return errs.MysqlOperateError
	}

	//获取所有外部id
	surplusOpenIds := []string{}
	outToInnerUserId := map[string]int64{}
	for _, info := range *outUserInfo {
		surplusOpenIds = append(surplusOpenIds, info.OutUserId)
		outToInnerUserId[info.OutUserId] = info.UserId
	}
	if len(surplusOpenIds) == 0 {
		return nil
	}

	scopeUsersReply, sdkErr := client.GetScopeUsers(&sdkVo.GetScopeUsersReq{Depts: outDeptInfos, OpenIds: surplusOpenIds})
	if sdkErr != nil {
		log.Errorf("[GetScopeUsers] err:%v", sdkErr)
		return errs.PlatFormOpenApiCallError
	}
	delRelationIds := []int64{}
	insertRelationPos := []po.PpmOrgUserDepartment{}
	for _, user := range scopeUsersReply.Users {
		userId, ok := outToInnerUserId[user.UserId]
		if !ok {
			continue
		}
		curDeptIds := []int64{}
		if user.DeptIds != nil {
			for _, department := range user.DeptIds {
				if deptId, ok := outToInnerDeptId[department]; ok {
					curDeptIds = append(curDeptIds, deptId)
				}
			}
		}
		// 去重
		curDeptIds = slice.SliceUniqueInt64(curDeptIds)
		beforeDeptIds := []int64{}
		if ids, ok := userDeptMap[userId]; ok {
			beforeDeptIds = ids
		}
		// 去重
		beforeDeptIds = slice.SliceUniqueInt64(beforeDeptIds)
		delDeptIds, addDeptIds := util.GetDifMemberIds(beforeDeptIds, curDeptIds)

		for _, id := range delDeptIds {
			delRelationIds = append(delRelationIds, userDeptRelationMap[userId][id]...)
		}
		for _, id := range addDeptIds {
			insertRelationPos = append(insertRelationPos, po.PpmOrgUserDepartment{
				Id:           0,
				OrgId:        orgId,
				UserId:       userId,
				DepartmentId: id,
			})
		}
	}

	if len(insertRelationPos) > 0 || len(delRelationIds) > 0 {
		if len(insertRelationPos) > 0 {
			userDepPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(insertRelationPos))
			if idErr != nil {
				log.Error(idErr)
				return idErr
			}
			for i, _ := range insertRelationPos {
				insertRelationPos[i].Id = userDepPoIds.Ids[i].Id
			}
		}

		err := mysql.TransX(func(tx sqlbuilder.Tx) error {
			if len(insertRelationPos) > 0 {
				err := PaginationInsert(slice.ToSlice(insertRelationPos), &po.PpmOrgUserDepartment{}, tx)
				if err != nil {
					log.Error(err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}

			if len(delRelationIds) > 0 {
				_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
					consts.TcOrgId:    orgId,
					consts.TcId:       db.In(delRelationIds),
					consts.TcIsDelete: consts.AppIsNoDelete,
				}, mysql.Upd{
					consts.TcIsDelete: consts.AppIsDeleted,
				})
				if err != nil {
					log.Errorf("[handleSyncUserDept] err: %v", err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}

			return nil
		})

		if err != nil {
			log.Errorf("[handleSyncUserDept] err: %v", err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	return nil
}
