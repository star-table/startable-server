package orgsvc

import (
	"strings"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"gitea.bjx.cloud/allstar/polaris-backend/service/platform/orgsvc/po"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/md5"
	"github.com/star-table/startable-server/common/core/util/pinyin"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/uservo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// InitThirdPlatformUsers 同步第三方平台用户到极星
func InitThirdPlatformUsers(client sdk_interface.Sdk, sourceChannel string, orgId int64, ignoredOpenIdMap map[string]int64, thirdDeptIdToPoId map[string]int64,
	depts []*sdkVo.DepartmentInfo, outOrgOwnerId string, hasSuperAdmin bool, tx sqlbuilder.Tx) errs.SystemErrorInfo {

	// 获取所有授权用户
	reply, err := client.GetScopeUsers(&sdkVo.GetScopeUsersReq{Depts: depts})
	if err != nil {
		log.Errorf("[initThirdPlatformUsers] GetScopeUsers err: %v", err)
		return errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	}
	log.Infof("[InitThirdPlatformUsers] 初始化第三方组织时，orgId: %d, installer openId: %s, 待始化用户的 openIds: %s", orgId, outOrgOwnerId, json.ToJsonIgnoreError(reply.Users))

	deptIdsMap := make(map[string]bool, len(depts))
	for _, dept := range depts {
		deptIdsMap[dept.Id] = true
	}

	if len(reply.Users) > 0 {
		openIds := make([]string, 0, len(reply.Users))
		usersMap := make(map[string]*sdkVo.ScopeUser, len(reply.Users))
		for _, user := range reply.Users {
			if _, ok := ignoredOpenIdMap[user.OpenId]; !ok {
				openIds = append(openIds, user.OpenId)
				usersMap[user.OpenId] = user
			}
		}
		if len(openIds) > 0 {
			filterOpenIds, err2 := FilterBoundingOpenIds(orgId, sourceChannel, openIds, tx)
			if err2 != nil {
				log.Errorf("[InitFsUserList] FilterBoundingOpenIds err: %v", err2)
				return err2
			}
			if len(filterOpenIds) > 0 {
				err2 = initPoUsers(sourceChannel, orgId, filterOpenIds, usersMap, thirdDeptIdToPoId, outOrgOwnerId, hasSuperAdmin, tx)
				if err2 != nil {
					log.Errorf("[InitFsUserList] initPoUsers err: %v", err2)
					return err2
				}
			}
		}
	}

	return nil
}

func initPoUsers(sourceChannel string, orgId int64, openIds []string, usersMap map[string]*sdkVo.ScopeUser, thirdDeptIdToPoId map[string]int64,
	outOrgOwnerId string, hasSuperAdmin bool, tx sqlbuilder.Tx) errs.SystemErrorInfo {

	poIdsMap, idErr := applyMultiplePrimaryIds([]string{consts.TableUser, consts.TableUserOutInfo,
		consts.TableUserConfig, consts.TableUserOrganization}, len(openIds))
	if idErr != nil {
		return idErr
	}

	userPoList := make([]po.PpmOrgUser, 0, len(openIds))
	userOutPoList := make([]po.PpmOrgUserOutInfo, 0, len(openIds))
	userConfigList := make([]po.PpmOrgUserConfig, 0, len(openIds))
	userOrgList := make([]po.PpmOrgUserOrganization, 0, len(openIds))
	userDepList := make([]po.PpmOrgUserDepartment, 0, len(openIds))
	defaultOrgOwnerId := int64(0)
	adminUserIds := make([]int64, 0, 3)

	for i, id := range openIds {
		poUserId := poIdsMap[consts.TableUser].Ids[i].Id
		scopeUser := usersMap[id]
		if scopeUser.OpenId == outOrgOwnerId {
			defaultOrgOwnerId = poUserId
		}
		if scopeUser.IsAdmin {
			adminUserIds = append(adminUserIds, poUserId)
		}

		userPoList = append(userPoList, assemblyOrgUserInfo(orgId, poUserId, sourceChannel, scopeUser))
		userOutPoList = append(userOutPoList, assemblyOrgOutInfo(orgId, poIdsMap[consts.TableUserOutInfo].Ids[i].Id, poUserId, sourceChannel, scopeUser))
		userConfigList = append(userConfigList, assemblyOrgUserConfigInfo(orgId, poIdsMap[consts.TableUserConfig].Ids[i].Id, poUserId))
		userOrgList = append(userOrgList, assemblyUserOrgRelationInfo(orgId, poIdsMap[consts.TableUserOrganization].Ids[i].Id, poUserId))

		for _, userDep := range scopeUser.DeptIds {
			if depId, ok := thirdDeptIdToPoId[userDep]; ok {
				//这些用户归入部门
				userDepList = append(userDepList, assemblyUserDepRelationInfo(orgId, poUserId, depId))
			}
		}
	}

	log.Infof("[initPoUsers] 初始化用户 管理员列表:%v, openIds:%v, orgId:%v", adminUserIds, openIds, orgId)

	// 如果有管理员，取第一个管理员为owner，如果没有管理员，则取授权的用户当管理员以及owner (hasSuperAdmin代表是否有owner，没有owner才创建)
	firstPoId := poIdsMap[consts.TableUser].Ids[0].Id // 如果安装人不在授权范围内的时候，用第一个人作为安装人
	err := initAdminAndOwner(orgId, defaultOrgOwnerId, firstPoId, adminUserIds, hasSuperAdmin, tx)
	if err != nil {
		return err
	}

	// 将数据存到db
	err = savePoUsersToDb(userPoList, userOutPoList, userOrgList, userConfigList, userDepList, tx)
	if err != nil {
		return err
	}

	return nil
}

// savePoUsersToDb 将授权用户写入到db中
func savePoUsersToDb(userPoList []po.PpmOrgUser, userOutPoList []po.PpmOrgUserOutInfo, userOrgList []po.PpmOrgUserOrganization,
	userConfigList []po.PpmOrgUserConfig, userDepList []po.PpmOrgUserDepartment, tx sqlbuilder.Tx) errs.SystemErrorInfo {

	err := PaginationInsert(slice.ToSlice(userPoList), &po.PpmOrgUser{}, tx)
	if err != nil {
		log.Errorf("[savePoUsersToDb] PaginationInsert err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	log.Infof("[savePoUsersToDb] userOutList:%v", json.ToJsonIgnoreError(userOutPoList))

	err = PaginationInsert(slice.ToSlice(userOutPoList), &po.PpmOrgUserOutInfo{}, tx)
	if err != nil {
		log.Errorf("[savePoUsersToDb] PaginationInsert err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = PaginationInsert(slice.ToSlice(userOrgList), &po.PpmOrgUserOrganization{}, tx)
	if err != nil {
		log.Errorf("[savePoUsersToDb] PaginationInsert err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = PaginationInsert(slice.ToSlice(userConfigList), &po.PpmOrgUserConfig{}, tx)
	if err != nil {
		log.Errorf("[savePoUsersToDb] PaginationInsert err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	// 由于一个用户属于多个部门，所以申请的id数量跟上面的不一样，需要单独处理
	userDepPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(userDepList))
	if idErr != nil {
		log.Errorf("[InitFsUserList] idfacade.ApplyMultiplePrimaryIdRelaxed err: %v", idErr)
		return idErr
	}
	for i := range userDepList {
		userDepList[i].Id = userDepPoIds.Ids[i].Id
	}
	err = PaginationInsert(slice.ToSlice(userDepList), &po.PpmOrgUserDepartment{}, tx)
	if err != nil {
		log.Errorf("[InitFsUserList] PaginationInsert err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func initAdminAndOwner(orgId, defaultOrgOwnerId, firstPoId int64, adminUserIds []int64, hasSuperAdmin bool, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	// 如果有管理员，取第一个管理员为owner，如果没有管理员，则取授权的用户当管理员以及owner (hasSuperAdmin代表是否有owner，没有owner才创建)
	if len(adminUserIds) > 0 {
		err := InitFsManager(orgId, adminUserIds)
		if err != nil {
			log.Errorf("[initPoUsers] InitFsManager, orgId:%v, adminUserIds:%v, err:%v", orgId, adminUserIds, err)
			return err
		}
		// 如果没有owner，让第一个管理员成为owner
		if !hasSuperAdmin {
			err = OrgOwnerInit(orgId, adminUserIds[0], defaultOrgOwnerId, tx)
			if err != nil {
				return err
			}
		}
	} else {
		if !hasSuperAdmin {
			// 不存在管理员，而且找不到安装人的时候，取第一个人
			if defaultOrgOwnerId == 0 {
				defaultOrgOwnerId = firstPoId
			}
			err := InitFsManager(orgId, []int64{defaultOrgOwnerId})
			if err != nil {
				log.Errorf("[initPoUsers] InitFsManager, orgId:%v, adminUserIds:%v, err:%v", orgId, defaultOrgOwnerId, err)
				return err
			}
			err = OrgOwnerInit(orgId, defaultOrgOwnerId, defaultOrgOwnerId, tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func applyMultiplePrimaryIds(keys []string, count int) (map[string]*bo.IdCodes, errs.SystemErrorInfo) {
	result := make(map[string]*bo.IdCodes, len(keys))
	for _, key := range keys {
		ids, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(key, count)
		if idErr != nil {
			log.Errorf("[applyMultiplePrimaryIds] idfacade.ApplyMultiplePrimaryIdRelaxed err: %v", idErr)
			return nil, idErr
		}
		result[key] = ids
	}

	return result, nil
}

// 分组插入新数据
func PaginationInsert(list []interface{}, domainObj mysql.Domain, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	totalSize := len(list)
	batch := 1000
	offset := 0
	for {
		limit := offset + batch
		if totalSize < limit {
			limit = totalSize
		}
		oneBatch := list[offset:limit]
		batchInsert := mysql.TransBatchInsert(tx, domainObj, oneBatch)
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}
		if totalSize <= limit {
			break
		}
		offset += batch
	}
	return nil
}

// 分组插入新数据
func PaginationInsertWithoutTrans(list []interface{}, domainObj mysql.Domain) errs.SystemErrorInfo {
	totalSize := len(list)
	batch := 1000
	offset := 0
	for {
		limit := offset + batch
		if totalSize < limit {
			limit = totalSize
		}
		oneBatch := list[offset:limit]
		batchInsert := mysql.BatchInsert(domainObj, oneBatch)
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}
		if totalSize <= limit {
			break
		}
		offset += batch
	}
	return nil
}

func InitFsManager(orgId int64, userIds []int64) errs.SystemErrorInfo {
	// 将 userId 对应的用户加入到管理组中，因为他是超管
	resp := userfacade.AddUserToSysManageGroup(orgId, userIds[0], uservo.AddUserToSysManageGroupReq{
		UserIds: userIds,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	return nil
}

// InitPlatformUser 我是真的不知道这个在搅什么
func InitPlatformUser(sourceChannel string, orgId int64, corpId, openUserId string, tx sqlbuilder.Tx, accessToken string, deptIds ...string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	client, err := platform_sdk.GetClient(sourceChannel, corpId)
	if err != nil {
		log.Errorf("[GetClient] err:%v", err)
		return nil, errs.PlatFormOpenApiCallError
	}

	checkReply, err := client.CheckIsAdminAuth()
	if err != nil {
		log.Errorf("[GetClient] CheckIsAdminAuth err:%v", err)
		return nil, errs.PlatFormOpenApiCallError
	}

	fsUser := &sdkVo.ScopeUser{}
	if checkReply.IsAdminAuth {
		var err2 errs.SystemErrorInfo
		fsUser, err2 = GetPlatformUserDetailInfo(client, openUserId, deptIds...)
		if err2 != nil {
			return nil, err2
		}
	} else {
		if accessToken == "" {
			//只有授权登录的时候能够拿到信息
			return &bo.UserInfoBo{}, nil
		}
		userInfo, err := sdk.GetOAuth2UserInfo(accessToken)
		if err != nil {
			log.Error(err)
			return nil, errs.FeiShuOpenApiCallError
		}
		fsUser = &sdkVo.ScopeUser{
			Name:   userInfo.Name,
			OpenId: openUserId,
			UserId: openUserId,
			Avatar: userInfo.AvatarUrl,
			Email:  userInfo.Email,
			Mobile: userInfo.Mobile,
		}
	}

	userNativeId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if idErr != nil {
		log.Error(idErr)
		return nil, idErr
	}
	userOutInfoId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
	if idErr != nil {
		log.Error(idErr)
		return nil, idErr
	}
	userConfigId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserConfig)
	if idErr != nil {
		log.Error(idErr)
		return nil, idErr
	}
	userOrgId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
	if idErr != nil {
		log.Error(idErr)
		return nil, idErr
	}
	insertUserDeptPos := []po.PpmOrgUserDepartment{}
	if fsUser.DeptIds != nil && len(fsUser.DeptIds) > 0 {
		//查找部门
		depts := &[]po.PpmOrgDepartmentOutInfo{}
		conn, err := mysql.GetConnect()
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		selectErr := conn.Select("o.department_id").From("ppm_org_department_out_info o", "ppm_org_department d").Where(db.Cond{
			"o." + consts.TcIsDelete:           consts.AppIsNoDelete,
			"d." + consts.TcIsDelete:           consts.AppIsNoDelete,
			"o." + consts.TcOrgId:              orgId,
			"d." + consts.TcOrgId:              orgId,
			"o." + consts.TcStatus:             consts.AppStatusEnable,
			"d." + consts.TcStatus:             consts.AppStatusEnable,
			"o." + consts.TcDepartmentId:       db.Raw("d." + consts.TcId),
			"o." + consts.TcOutOrgDepartmentId: db.In(fsUser.DeptIds),
		}).All(depts)
		if selectErr != nil {
			log.Error(selectErr)
			return nil, errs.MysqlOperateError
		}

		if len(*depts) > 0 {
			userDepIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(*depts))
			if idErr != nil {
				log.Error(idErr)
				return nil, idErr
			}
			for i, info := range *depts {
				insertUserDeptPos = append(insertUserDeptPos, po.PpmOrgUserDepartment{
					Id:           userDepIds.Ids[i].Id,
					OrgId:        orgId,
					UserId:       userNativeId,
					DepartmentId: info.DepartmentId,
				})
			}
		}
	}

	userPo := assemblyOrgUserInfo(orgId, userNativeId, sourceChannel, fsUser)
	userOutPo := assemblyOrgOutInfo(orgId, userOutInfoId, userNativeId, sourceChannel, fsUser)
	userConfigPo := assemblyOrgUserConfigInfo(orgId, userConfigId, userNativeId)
	userOrgRelationPo := assemblyUserOrgRelationInfo(orgId, userOrgId, userNativeId)

	dbErr := mysql.TransInsert(tx, &userPo)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userOutPo)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userConfigPo)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userOrgRelationPo)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	if len(insertUserDeptPos) > 0 {
		dbErr = mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, slice.ToSlice(insertUserDeptPos))
		if dbErr != nil {
			log.Error(dbErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}

	userBo := &bo.UserInfoBo{}
	_ = copyer.Copy(userPo, userBo)
	return userBo, nil
}

func AssemblyFeiShuUserInfo(orgId, poId int64, scopeUser *sdkVo.ScopeUser) po.PpmOrgUser {
	return assemblyOrgUserInfo(orgId, poId, sdk_const.SourceChannelFeishu, scopeUser)
}

func assemblyOrgUserInfo(orgId, poId int64, sourceChannel string, scopeUser *sdkVo.ScopeUser) po.PpmOrgUser {
	phoneNumber := scopeUser.Mobile
	sourcePlatform := sourceChannel
	name := scopeUser.Name

	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userPo := &po.PpmOrgUser{
		Id:                 poId,
		OrgId:              orgId,
		Name:               name,
		NamePinyin:         pinyin.ConvertToPinyin(name),
		Avatar:             scopeUser.Avatar,
		LoginName:          phoneNumber, //
		LoginNameEditCount: 0,
		Email:              scopeUser.Email,
		Mobile:             phoneNumber,
		Password:           pwd,
		PasswordSalt:       salt,
		SourceChannel:      sourceChannel,
		SourcePlatform:     sourcePlatform,
	}
	return *userPo
}

func assemblyOrgOutInfo(orgId, poId, userId int64, sourceChannel string, scopeUser *sdkVo.ScopeUser) po.PpmOrgUserOutInfo {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOutInfo := &po.PpmOrgUserOutInfo{}
	userOutInfo.Id = poId
	userOutInfo.UserId = userId
	userOutInfo.OrgId = orgId
	// 之前飞书两个都是openId，所以飞书情况下两个都是都是openId，钉钉应该大部分用的是userId，而不是openId（unionId）
	userOutInfo.OutOrgUserId = scopeUser.OpenId
	userOutInfo.OutUserId = scopeUser.UserId
	userOutInfo.IsDelete = consts.AppIsNoDelete
	userOutInfo.Status = consts.AppStatusEnable
	userOutInfo.SourceChannel = sourceChannel
	userOutInfo.Name = scopeUser.Name
	userOutInfo.Avatar = scopeUser.Avatar
	userOutInfo.JobNumber = scopeUser.JobNumber

	return *userOutInfo
}

func assemblyOrgUserConfigInfo(orgId, poId, userId int64) po.PpmOrgUserConfig {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userConfigInfo := &po.PpmOrgUserConfig{}
	userConfigInfo.Id = poId
	userConfigInfo.UserId = userId
	userConfigInfo.OrgId = orgId

	return *userConfigInfo
}

func assemblyUserOrgRelationInfo(orgId, poId, userId int64) po.PpmOrgUserOrganization {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOrgRelationInfo := &po.PpmOrgUserOrganization{}
	userOrgRelationInfo.Id = poId
	userOrgRelationInfo.UserId = userId
	userOrgRelationInfo.OrgId = orgId
	userOrgRelationInfo.Status = consts.AppStatusEnable
	userOrgRelationInfo.UseStatus = consts.AppStatusDisabled
	userOrgRelationInfo.CheckStatus = consts.AppCheckStatusSuccess

	return *userOrgRelationInfo
}

func assemblyUserDepRelationInfo(orgId int64, userId int64, depId int64) po.PpmOrgUserDepartment {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userDepRelationInfo := &po.PpmOrgUserDepartment{}
	userDepRelationInfo.UserId = userId
	userDepRelationInfo.OrgId = orgId
	userDepRelationInfo.DepartmentId = depId

	return *userDepRelationInfo
}

func InitInstaller(orgId int64, outOrgOwnerId, sourceChannel string, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//用户表
	userId, userIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if userIdErr != nil {
		log.Error(userIdErr)
		return userIdErr
	}
	//外部用户表
	outUserId, outUserIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
	if outUserIdErr != nil {
		log.Error(outUserIdErr)
		return outUserIdErr
	}
	//用户配置表
	userConfigId, userConfigIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserConfig)
	if userConfigIdErr != nil {
		log.Error(userConfigIdErr)
		return userConfigIdErr
	}
	//组织用户表
	userOrgId, userOrgIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
	if userOrgIdErr != nil {
		log.Error(userOrgIdErr)
		return userOrgIdErr
	}

	//组织用户表
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err1 := mysql.TransInsert(tx, &po.PpmOrgUser{
			Id:             userId,
			OrgId:          orgId,
			Name:           "未激活",
			NamePinyin:     pinyin.ConvertToPinyin("未激活"),
			Avatar:         consts.AvatarForUnallocated,
			SourcePlatform: "",
			SourceChannel:  sourceChannel,
		})
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		err2 := mysql.TransInsert(tx, &po.PpmOrgUserOutInfo{
			Id:             outUserId,
			OrgId:          orgId,
			UserId:         userId,
			SourcePlatform: "",
			SourceChannel:  sourceChannel,
			OutOrgUserId:   outOrgOwnerId,
			OutUserId:      outOrgOwnerId,
			Name:           "未激活",
			Avatar:         consts.AvatarForUnallocated,
		})
		if err2 != nil {
			log.Error(err2)
			return err2
		}

		err3 := mysql.TransInsert(tx, &po.PpmOrgUserConfig{
			Id:     userConfigId,
			OrgId:  orgId,
			UserId: userId,
		})
		if err3 != nil {
			log.Error(err3)
			return err3
		}

		err4 := mysql.TransInsert(tx, &po.PpmOrgUserOrganization{
			Id:          userOrgId,
			OrgId:       orgId,
			UserId:      userId,
			CheckStatus: consts.AppCheckStatusSuccess,
			UseStatus:   consts.AppStatusDisabled,
			Status:      consts.AppStatusEnable,
		})
		if err4 != nil {
			log.Error(err4)
			return err4
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return errs.MysqlOperateError
	}

	err2 := InitFsManager(orgId, []int64{userId})
	if err2 != nil {
		log.Error(err2)
		return err2
	}

	err2 = OrgOwnerInit(orgId, userId, userId, tx)
	if err2 != nil {
		log.Error(err2)
		return err2
	}

	return nil
}

func ScheduleOrgUserMobileAndEmail(tenantKey string) errs.SystemErrorInfo {
	tenant, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}

	orgInfo, err := GetOrgInfoByOutOrgId(tenantKey, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Error(err)
		return err
	}

	surplusOpenIds, err := feishu.GetScopeOpenIdsLimit(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}

	batch := 100
	offset := 0
	surplusSize := len(surplusOpenIds)

	//默认传过来的组织拥有者id
	if surplusSize > 0 {
		for {
			limit := offset + batch
			if surplusSize < limit {
				limit = surplusSize
			}
			openIds := surplusOpenIds[offset:limit]

			userBatchResp, err := tenant.GetUserBatchGetV2(nil, openIds)
			if err != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if userBatchResp.Code != 0 {
				log.Error(userBatchResp.Msg)
				return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}

			for _, fsUser := range userBatchResp.Data.Users {
				uid, err := GetUserIdByEmpId(orgInfo.OrgId, fsUser.OpenId)
				if err != nil {
					log.Error(err)
					return err
				}
				fsUser.Mobile = strings.ReplaceAll(fsUser.Mobile, "+86", "")
				err = UpdateUserInfo(uid, mysql.Upd{
					consts.TcMobile: fsUser.Mobile,
					consts.TcEmail:  fsUser.Email,
				})
				if err != nil {
					log.Error(err)
					return err
				}
			}

			if surplusSize <= limit {
				break
			}
			offset += batch
		}
	}
	return nil
}

func fsUserToScopeUser(user vo.UserDetailInfoV2) *sdkVo.ScopeUser {
	return &sdkVo.ScopeUser{
		OpenId:    user.OpenId,
		UserId:    user.UserId,
		Name:      user.Name,
		JobNumber: user.EmployeeNo,
		Avatar:    user.Avatar.AvatarOrigin,
		Email:     user.Email,
		Mobile:    user.Mobile,
		DeptIds:   user.Departments,
		IsAdmin:   user.IsTenantManager,
	}
}

// GetPlatformUserDetailInfo 钉钉outUserId存的是userId，而飞书存的是openId
func GetPlatformUserDetailInfo(client sdk_interface.Sdk, outUserId string, deptIds ...string) (*sdkVo.ScopeUser, errs.SystemErrorInfo) {
	ids := []string{outUserId}
	deptId := "0"
	if len(deptIds) > 0 {
		deptId = deptIds[0]
	}
	userInfos, err := client.GetUserDetailInfos(&sdkVo.GetUserDetailInfosReq{Ids: ids, DeptId: deptId})
	if err != nil {
		log.Errorf("[GetPlatformUserDetailInfo] GetUserDetailInfos ids:%v, err:%v", ids, err)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	}

	if len(userInfos.Users) == 0 {
		log.Errorf("[GetPlatformUserDetailInfo] user不存在")
		return nil, errs.FeiShuUserNotInAppUseScopeOfAuthority
	}
	if len(userInfos.Users[0].DeptIds) == 0 {
		userInfos.Users[0].DeptIds = deptIds
	}

	return userInfos.Users[0], nil
}
