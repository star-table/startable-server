package orgsvc

import (
	"errors"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/service/orgsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// InitDepartment 初始化部门
func InitDepartment(client sdk_interface.Sdk, orgId int64, corpId string, sourceChannel string,
	outOrgOwnerId string, tx sqlbuilder.Tx) errs.SystemErrorInfo {

	//首先判断是否有相应权限（没有权限表示是普通用户安装，则不初始化人员部门，auth的时候主动加入）
	checkReply, sdkErr := client.CheckIsAdminAuth()
	if sdkErr != nil {
		log.Errorf("[initDepartment] CheckIsAdminAuth err: %v", sdkErr)
		return errs.PlatFormOpenApiCallError
	}

	log.Infof("初始化安装 orgId: %d, 安装安装者 openId: %s, tenantKey: %s, isInScope: %v", orgId, outOrgOwnerId, corpId, checkReply)
	if !checkReply.IsAdminAuth && sourceChannel == sdk_const.SourceChannelFeishu {
		//把安装者作为唯一的用户，因为钉钉平台没有安装者这个逻辑，所以在拿授权信息的时候取第一个授权的人
		if outOrgOwnerId == "" {
			outOrgOwnerId = checkReply.FirstAuthUserId
		}
		if outOrgOwnerId != "" {
			err := InitInstaller(orgId, outOrgOwnerId, sourceChannel, tx)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		return nil
	}

	deptList, sdkErr := client.GetScopeDeps()
	if sdkErr != nil {
		log.Errorf("[initDepartment] client.GetScopeDeps err: %v", sdkErr)
		return errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}
	deptMap := map[string]*sdkVo.DepartmentInfo{}
	for _, dep := range deptList.Depts {
		deptMap[dep.Id] = dep
	}

	depSize := len(deptList.Depts)
	departmentInfo := make([]interface{}, depSize)
	outDepartmentInfo := make([]interface{}, depSize)
	thirdDeptIdToPoId := map[string]int64{}

	depIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartment, depSize)
	if err != nil {
		log.Error(err)
		return err
	}

	depOutIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartmentOutInfo, depSize)
	if err != nil {
		log.Error(err)
		return err
	}

	for k, v := range deptList.Depts {
		thirdDeptIdToPoId[v.Id] = depIds.Ids[k].Id
	}

	rootId := int64(0)
	for k, v := range deptList.Depts {
		depId := depIds.Ids[k].Id
		depOutId := depOutIds.Ids[k].Id
		parentDepId := int64(0)

		if id, ok := thirdDeptIdToPoId[v.ParentId]; ok {
			parentDepId = id
		} else {
			parentDepId = rootId
		}
		if depId == rootId {
			parentDepId = 0
		}

		departmentInfo[k] = &po.PpmOrgDepartment{
			Id:            depId,
			OrgId:         orgId,
			Name:          v.Name,
			ParentId:      parentDepId,
			SourceChannel: sourceChannel,
		}

		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
			Id:                       depOutId,
			OrgId:                    orgId,
			DepartmentId:             depId,
			SourceChannel:            sourceChannel,
			OutOrgDepartmentId:       v.Id,
			Name:                     v.Name,
			OutOrgDepartmentParentId: v.ParentId,
		}
	}

	//初始化用户
	userInitErr := InitThirdPlatformUsers(client, sourceChannel, orgId, map[string]int64{}, thirdDeptIdToPoId,
		deptList.Depts, outOrgOwnerId, false, tx)
	if userInitErr != nil {
		log.Error(userInitErr)
		return userInitErr
	}

	departErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
	if departErr != nil {
		log.Error(departErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
	}
	outDepartErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
	if outDepartErr != nil {
		log.Error(outDepartErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
	}

	return nil
}

func BoundFsDepartment(orgId int64, tenantKey string, tx sqlbuilder.Tx) (map[string]int64, errs.SystemErrorInfo) {
	// 清理原来的组织部门
	_, mysqlErr := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
		consts.TcOrgId: orgId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}
	_, mysqlErr = mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{
		consts.TcOrgId: orgId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}
	// 清理原来的部门用户关联
	_, mysqlErr = mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
		consts.TcOrgId: orgId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}

	// 绑定飞书组织
	deptList, err := feishu.GetScopeDeps(tenantKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	deptMap := map[string]vo.DepartmentRestInfoVo{}
	for _, dep := range deptList {
		deptMap[dep.Id] = dep
	}
	depSize := len(deptList)
	departmentInfo := make([]interface{}, len(deptList))
	outDepartmentInfo := make([]interface{}, len(deptList))
	fsDepIdMap := map[string]int64{}
	depIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartment, depSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	depOutIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartmentOutInfo, depSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for k, v := range deptList {
		fsDepIdMap[v.Id] = depIds.Ids[k].Id
	}
	rootId := int64(0)
	for k, v := range deptList {
		depId := depIds.Ids[k].Id
		depOutId := depOutIds.Ids[k].Id
		parentDepId := int64(0)
		if id, ok := fsDepIdMap[v.ParentId]; ok {
			parentDepId = id
		} else {
			parentDepId = rootId
		}
		if depId == rootId {
			parentDepId = 0
		}
		departmentInfo[k] = &po.PpmOrgDepartment{
			Id:            depId,
			OrgId:         orgId,
			Name:          v.Name,
			ParentId:      parentDepId,
			SourceChannel: sdk_const.SourceChannelFeishu,
		}
		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
			Id:                       depOutId,
			OrgId:                    orgId,
			DepartmentId:             depId,
			SourceChannel:            sdk_const.SourceChannelFeishu,
			OutOrgDepartmentId:       v.Id,
			Name:                     v.Name,
			OutOrgDepartmentParentId: v.ParentId,
		}
	}
	departErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
	if departErr != nil {
		log.Error(departErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
	}
	outDepartErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
	if outDepartErr != nil {
		log.Error(outDepartErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
	}
	return fsDepIdMap, nil
}

func BoundFsUser(orgId int64, creator int64, tenantKey string, deptMapping map[string]int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	// 绑定飞书用户
	outUserMapping, err := LocalUserBoundOutInfo(orgId, creator, tenantKey, deptMapping, tx)
	if err != nil {
		log.Error(err)
		return err
	}
	// 同步其他fs用户
	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelFeishu, tenantKey)
	if sdkErr != nil {
		log.Errorf("[BoundFsUser] platform_sdk.GetClient err: %v", sdkErr)
		return errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}

	userInitErr := InitThirdPlatformUsers(client, sdk_const.SourceChannelFeishu, orgId, outUserMapping, deptMapping, nil, "", true, tx)
	if userInitErr != nil {
		log.Error(userInitErr)
		return userInitErr
	}
	return nil
}

// LocalUserBoundOutInfo 本地用户绑定飞书
// 通过本地用户的手机号邮箱绑定飞书用户
// 未匹配到的用户隐藏
func LocalUserBoundOutInfo(orgId int64, creator int64, tenantKey string, deptMapping map[string]int64, tx sqlbuilder.Tx) (map[string]int64, errs.SystemErrorInfo) {
	// 绑定飞书用户
	page := 1
	batch := 50
	missingUserIds := make([]int64, 0)
	outUserIds := map[string]int64{}
	for {
		userOrgs := make([]po.PpmOrgUserOrganization, 0)
		_, mysqlErr := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		}, nil, page, batch, nil, &userOrgs)
		if mysqlErr != nil {
			log.Error(mysqlErr)
			return nil, errs.MysqlOperateError
		}
		userIds := make([]int64, 0)
		for _, userOrg := range userOrgs {
			if userOrg.UserId == creator {
				continue
			}
			userIds = append(userIds, userOrg.UserId)
		}
		if len(userIds) > 0 {
			users := make([]po.PpmOrgUser, 0)
			mysqlErr = mysql.SelectAllByCond(consts.TableUser, db.Cond{
				consts.TcId:       db.In(userIds),
				consts.TcIsDelete: consts.AppIsNoDelete,
			}, &users)
			if mysqlErr != nil {
				log.Error(mysqlErr)
				return nil, errs.MysqlOperateError
			}
			// 获取飞书的用户信息
			emails := make([]string, 0)
			mobiles := make([]string, 0)
			mapping := map[string]po.PpmOrgUser{}
			for _, user := range users {
				if user.Email != "" {
					emails = append(emails, user.Email)
					mapping[user.Email] = user
				}
				if user.Mobile != "" {
					mobile := util.GetMobile(user.Mobile)
					mobiles = append(mobiles, mobile)
					mapping[mobile] = user
				}
			}
			tenant, err := feishu.GetTenant(tenantKey)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			userMapping := map[int64]string{}
			openIds := make([]string, 0)
			if len(emails) > 0 || len(mobiles) > 0 {
				batchGetIdResp, fsErr := tenant.BatchGetId(emails, mobiles)
				if fsErr != nil {
					log.Error(fsErr)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, fsErr)
				}
				if batchGetIdResp.Code != 0 {
					log.Error(batchGetIdResp.Msg)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, errors.New(batchGetIdResp.Msg))
				}
				for email, fsUserId := range batchGetIdResp.Data.EmailUsers {
					userMapping[mapping[email].Id] = fsUserId[0].OpenId
					openIds = append(openIds, fsUserId[0].OpenId)
				}
				for mobile, fsUserId := range batchGetIdResp.Data.MobileUsers {
					userMapping[mapping[mobile].Id] = fsUserId[0].OpenId
					openIds = append(openIds, fsUserId[0].OpenId)
				}
			}
			openIds, err = FilterBoundingOpenIds(orgId, sdk_const.SourceChannelFeishu, openIds, tx)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			openIdSet := map[string]bool{}
			for _, openId := range openIds {
				openIdSet[openId] = true
			}
			for _, userId := range userIds {
				if openId, ok := userMapping[userId]; !ok { // 未匹配到，隐藏本地用户
					missingUserIds = append(missingUserIds, userId)
				} else {
					outUserIds[openId] = userId //匹配到了，增加openId和本地userId的映射
				}
			}
			if len(openIds) > 0 {
				userBatchResp, fsErr := tenant.GetUserBatchGetV2(nil, openIds)
				if fsErr != nil {
					log.Error(fsErr)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, fsErr)
				}
				if userBatchResp.Code != 0 {
					log.Error(userBatchResp.Msg)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, errors.New(userBatchResp.Msg))
				}
				fsUsers := map[string]*sdkVo.ScopeUser{}
				for _, fsUser := range userBatchResp.Data.Users {
					fsUsers[fsUser.OpenId] = fsUserToScopeUser(fsUser)
				}

				// 为用户封装外部信息
				userOutInfos := make([]po.PpmOrgUserOutInfo, 0)
				userDeptInfos := make([]po.PpmOrgUserDepartment, 0)
				for userId, fsOpenId := range userMapping {
					if openIdSet[fsOpenId] {
						if fsUser, ok := fsUsers[fsOpenId]; ok {
							outUserInfo := assemblyOrgOutInfo(orgId, 0, userId, sdk_const.SourceChannelFeishu, fsUser)
							userOutInfos = append(userOutInfos, outUserInfo)
							for _, fsDeptId := range fsUser.DeptIds {
								if deptId, ok := deptMapping[fsDeptId]; ok {
									userDeptInfo := assemblyUserDepRelationInfo(orgId, userId, deptId)
									userDeptInfos = append(userDeptInfos, userDeptInfo)
								}
							}
						}
					}
				}
				if len(userOutInfos) > 0 {
					userOutPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOutInfo, len(userOutInfos))
					if idErr != nil {
						log.Error(idErr)
						return nil, idErr
					}
					for i, _ := range userOutInfos {
						userOutInfos[i].Id = userOutPoIds.Ids[i].Id
					}
					err = PaginationInsert(slice.ToSlice(userOutInfos), &po.PpmOrgUserOutInfo{}, tx)
					if err != nil {
						log.Error(err)
						return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
					}
				}
				if len(userDeptInfos) > 0 {
					userDeptIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(userDeptInfos))
					if idErr != nil {
						log.Error(idErr)
						return nil, idErr
					}
					for i, _ := range userDeptInfos {
						userDeptInfos[i].Id = userDeptIds.Ids[i].Id
					}
					err = PaginationInsert(slice.ToSlice(userDeptInfos), &po.PpmOrgUserDepartment{}, tx)
					if err != nil {
						log.Error(err)
						return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
					}
				}
			}
		}
		if len(userOrgs) < batch {
			break
		}
		page++
	}
	log.Infof("组织 %d missionUserIds %v", orgId, json.ToJsonIgnoreError(missingUserIds))
	// 处理未匹配到的用户，隐藏，排除创建人
	if len(missingUserIds) > 0 {
		_, mysqlErr := mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
			consts.TcUserId:       db.In(missingUserIds),
			consts.TcUserId + " ": db.NotEq(creator),
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcStatus: consts.AppStatusHidden,
		})
		if mysqlErr != nil {
			log.Error(mysqlErr)
			return nil, errs.MysqlOperateError
		}
	}

	return outUserIds, nil
}

func BoundFsOutInfo(userId, orgId int64, name, openId string) errs.SystemErrorInfo {
	userOutInfoId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
	if err != nil {
		log.Error(err)
		return err
	}
	outUserInfo := assemblyOrgOutInfo(orgId, userOutInfoId, userId, sdk_const.SourceChannelFeishu, &sdkVo.ScopeUser{
		OpenId: openId,
		Name:   name,
	})
	outUserInfo.Id = userOutInfoId
	mysqlErr := mysql.Insert(&outUserInfo)
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return errs.MysqlOperateError
	}
	return nil
}

func BoundFsOutInfoTx(tx sqlbuilder.Tx, userId, orgId int64, name, openId string) errs.SystemErrorInfo {
	userOutInfoId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
	if err != nil {
		log.Error(err)
		return err
	}
	outUserInfo := assemblyOrgOutInfo(orgId, userOutInfoId, userId, sdk_const.SourceChannelFeishu, &sdkVo.ScopeUser{
		OpenId: openId,
		Name:   name,
	})
	outUserInfo.Id = userOutInfoId
	mysqlErr := mysql.TransInsert(tx, &outUserInfo)
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return errs.MysqlOperateError
	}
	return nil
}

func BoundFsAccount(userId, orgId int64, tenantKey, openId string) errs.SystemErrorInfo {
	tenant, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	userBatchResp, fsErr := tenant.GetUserBatchGetV2(nil, []string{openId})
	if fsErr != nil {
		log.Error(fsErr)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, fsErr)
	}
	if userBatchResp.Code != 0 {
		log.Error(userBatchResp.Msg)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, errors.New(userBatchResp.Msg))
	}
	if len(userBatchResp.Data.Users) == 0 {
		return errs.UserOutInfoNotExist
	}
	fsUser := userBatchResp.Data.Users[0]

	userOutInfoId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
	if err != nil {
		log.Error(err)
		return err
	}
	outUserInfo := assemblyOrgOutInfo(orgId, userOutInfoId, userId, sdk_const.SourceChannelFeishu, fsUserToScopeUser(fsUser))
	outUserInfo.Id = userOutInfoId
	mysqlErr := mysql.Insert(&outUserInfo)
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return errs.MysqlOperateError
	}

	// 找部门
	if len(fsUser.Departments) > 0 {
		deptIds, err := GetDepartmentIdsByOutDepartmentIds(orgId, fsUser.Departments)
		if err != nil {
			log.Error(err)
			return err
		}
		userDeptInfos := make([]po.PpmOrgUserDepartment, 0)
		for _, deptId := range deptIds {
			userDeptInfo := assemblyUserDepRelationInfo(orgId, userId, deptId)
			userDeptInfos = append(userDeptInfos, userDeptInfo)
		}
		if len(userDeptInfos) > 0 {
			userDeptIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(userDeptInfos))
			if idErr != nil {
				log.Error(idErr)
				return idErr
			}
			for i, _ := range userDeptInfos {
				userDeptInfos[i].Id = userDeptIds.Ids[i].Id
			}

			mysqlErr = mysql.BatchInsert(&po.PpmOrgUserDepartment{}, slice.ToSlice(userDeptInfos))
			if mysqlErr != nil {
				log.Error(mysqlErr)
				return errs.MysqlOperateError
			}
		}
	}
	return nil
}

// FilterBoundingOpenIds 过滤已被绑定的openId
func FilterBoundingOpenIds(orgId int64, sourceChannel string, openIds []string, tx sqlbuilder.Tx) ([]string, errs.SystemErrorInfo) {
	if len(openIds) == 0 {
		return openIds, nil
	}
	outInfos := make([]po.PpmOrgUserOutInfo, 0)
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: sourceChannel,
		consts.TcOutUserId:     db.In(openIds),
	}, &outInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	openIdSet := map[string]bool{}
	for _, v := range openIds {
		openIdSet[v] = true
	}
	filteredOutUser := make([]po.PpmOrgUserOutInfo, 0)
	alreadyBoundOrgUserId := map[int64]bool{}
	for _, outInfo := range outInfos {
		if _, ok := openIdSet[outInfo.OutUserId]; ok {
			if outInfo.OrgId == 0 {
				filteredOutUser = append(filteredOutUser, outInfo)
			} else {
				alreadyBoundOrgUserId[outInfo.UserId] = true
			}
			delete(openIdSet, outInfo.OutUserId)
		}
	}

	orphanOpenIds := make([]string, 0)
	for v := range openIdSet {
		orphanOpenIds = append(orphanOpenIds, v)
	}
	// 对于过滤的openId, 查出没有组织id的外部信息对应的用户并判断是否属于当前本地企业，属于的话为之创建一个新的外部信息
	if len(filteredOutUser) > 0 {
		userIds := make([]int64, 0)
		outUserMap := map[int64]po.PpmOrgUserOutInfo{}
		for _, outUser := range filteredOutUser {
			if alreadyBoundOrgUserId[outUser.UserId] {
				continue
			}
			userIds = append(userIds, outUser.UserId)
			outUserMap[outUser.UserId] = outUser
		}
		// 查询用户组织关联
		userOrgs := make([]po.PpmOrgUserOrganization, 0)
		err = mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
			consts.TcUserId:   db.In(userIds),
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &userOrgs)
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
		existOutUsers := make([]po.PpmOrgUserOutInfo, 0)
		for _, userOrg := range userOrgs {
			existOutUsers = append(existOutUsers, outUserMap[userOrg.UserId])
		}
		if len(existOutUsers) > 0 {
			userOutPoList := make([]po.PpmOrgUserOutInfo, 0)
			for _, existOutUser := range existOutUsers {
				outUserInfo := assemblyOrgOutInfo(orgId, 0, existOutUser.UserId, sourceChannel,
					&sdkVo.ScopeUser{
						OpenId: existOutUser.OutUserId,
						Name:   existOutUser.Name,
					})
				userOutPoList = append(userOutPoList, outUserInfo)
			}

			userOutPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOutInfo, len(userOutPoList))
			if idErr != nil {
				log.Error(idErr)
				return nil, idErr
			}
			for i, _ := range userOutPoList {
				userOutPoList[i].Id = userOutPoIds.Ids[i].Id
			}
			if tx != nil {
				insertErr := PaginationInsert(slice.ToSlice(userOutPoList), &po.PpmOrgUserOutInfo{}, tx)
				if insertErr != nil {
					log.Error(err)
					return nil, insertErr
				}
			} else {
				insertErr := PaginationInsertWithoutTrans(slice.ToSlice(userOutPoList), &po.PpmOrgUserOutInfo{})
				if insertErr != nil {
					log.Error(err)
					return nil, insertErr
				}
			}

		}
	}
	return orphanOpenIds, nil
}
