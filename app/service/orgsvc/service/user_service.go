package orgsvc

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sdkVo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/google/martian/log"
	"github.com/nyaruka/phonenumbers"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	sconsts "github.com/star-table/startable-server/app/service"
	service "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/app/service/orgsvc/po"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/md5"
	"github.com/star-table/startable-server/common/core/util/pinyin"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var userErrorTemplate = " GetCurrentUser : %v\n"

func PersonalInfo(orgId, userId int64, sourceChannel string) (*orgvo.PersonalInfo, errs.SystemErrorInfo) {
	orgConfig, err := GetOrgConfig(orgId)
	if err != nil {
		log.Errorf("[PersonalInfo] GetOrgConfig err: %v, orgId: %d", err, orgId)
		return nil, err
	}
	userInfoBo, passwordSet, err1 := domain.GetUserInfo(orgId, userId, sourceChannel)
	if err1 != nil {
		log.Errorf("[PersonalInfo] GetUserInfo err: %v, userId: %d", err1, userId)
		return nil, errs.BuildSystemErrorInfo(errs.GetUserInfoError, err1)
	}
	userInfoBo.Level = orgConfig.PayLevel
	userInfoBo.LevelName = businees.GetOrgPayVersionName(orgConfig.PayLevel)

	personalInfo := &orgvo.PersonalInfo{}
	copyErr := copyer.Copy(userInfoBo, personalInfo)
	if copyErr != nil {
		log.Errorf("[PersonalInfo] Copy err: %v", strs.ObjectToString(copyErr))
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	if passwordSet {
		personalInfo.PasswordSet = 1
	}

	if orgId > 0 {
		// 判断当前用户是不是超管
		adminFlagBo, err := service.GetUserAdminFlag(orgId, userId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		personalInfo.IsAdmin = adminFlagBo.IsAdmin
		personalInfo.IsManager = adminFlagBo.IsManager
		// 检查是否是平台管理员
		if ok, errSlice := slice.Contain([]string{sdk_const.SourceChannelFeishu, sdk_const.SourceChannelDingTalk,
			sdk_const.SourceChannelWeixin}, sourceChannel); errSlice == nil && ok {
			baseOrgInfo, err := GetBaseOrgInfo(orgId)
			if err != nil {
				log.Errorf("[PersonalInfo] GetBaseOrgInfo err: %v", err)
				return nil, err
			}
			if baseOrgInfo.OutOrgId != "" {
				platformAdmin, err := domain.GetPlatformAdmin(orgId, userId, baseOrgInfo.OutOrgId, *userInfoBo.EmplID, sourceChannel)
				if err != nil {
					// 如果不是该组织未安装应用，则忽略该错误。从而当前用户不是飞书的应用管理员
					log.Infof("[PersonalInfo] domain.GetFsPlatformAdmin err: %v, orgId: %d", err, orgId)
				}
				personalInfo.IsPlatformAdmin = platformAdmin
			}
		}
	}

	// 查询用户是否观看新手指南
	personalInfo.ExtraDataMap = make(map[string]interface{}, 0)
	userConfigBo, err := domain.GetUserConfigInfo(orgId, userId)
	if err != nil {
		if strings.Contains(err.Error(), db.ErrNoMoreRows.Error()) {
			// 没有则新增 user config 记录
			userConfigBo, err = domain.InsertUserConfig(orgId, userId)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		} else {
			log.Error(err)
			return nil, err
		}
	}
	if err := domain.GetNewUserGuideInfo(orgId, *userConfigBo, personalInfo.ExtraDataMap); err != nil {
		log.Errorf("[PersonalInfo] GetNewUserGuideInfo err: %v", err)
	}

	// 临时活动 双11 需要展示弹窗
	err = domain.GetActivity20221111Info(orgId, *userConfigBo, orgConfig.IsPayActivity11, personalInfo.ExtraDataMap)
	if err != nil {
		log.Errorf("[PersonalInfo] GetActivity20221111Info err: %v", err)
	}
	remindPopUp := domain.GetRemindPopUp(userConfigBo, personalInfo.RemindBindPhone)
	personalInfo.RemindPopUp = remindPopUp

	return personalInfo, nil
}

func GetUserIds(orgId int64, corpId, sourceChannel string, empIds []string) ([]*vo.UserIDInfo, errs.SystemErrorInfo) {
	resultIds := make([]*vo.UserIDInfo, len(empIds))
	for i, empId := range empIds {
		baseUserInfo, err := domain.GetBaseUserInfoByEmpId(orgId, empId)
		if err != nil {
			baseUserInfo, err = domain.UserInit(orgId, corpId, empId, sourceChannel)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.UserNotInitError, err)
			}
		}
		resultIds[i] = &vo.UserIDInfo{
			UserID:     baseUserInfo.UserId,
			Name:       baseUserInfo.Name,
			Avatar:     baseUserInfo.Avatar,
			EmplID:     baseUserInfo.OutUserId,
			IsDeleted:  baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted,
			IsDisabled: baseUserInfo.OrgUserStatus == consts.AppStatusDisabled,
		}
	}
	return resultIds, nil
}

func GetUserId(orgId int64, corpId, sourceChannel, empId string) (*vo.UserIDInfo, errs.SystemErrorInfo) {
	userIdInfos, err := GetUserIds(orgId, corpId, sourceChannel, []string{empId})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if userIdInfos == nil || len(userIdInfos) == 0 {
		log.Errorf("GetUserIds 获取到空 %d %s %s %s", orgId, corpId, sourceChannel, empId)
		return nil, errs.BuildSystemErrorInfo(errs.UserNotExist)
	}
	return userIdInfos[0], nil
}

func UserConfigInfo(orgId, userId int64) (*vo.UserConfig, errs.SystemErrorInfo) {
	//cacheUserInfo, err := GetCurrentUser(ctx)
	//if err != nil {
	//	log.Errorf(userErrorTemplate, err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId
	//userId := cacheUserInfo.UserId

	userConfig, err := domain.GetUserConfigInfo(orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	configInfo := &vo.UserConfig{}
	copyErr := copyer.Copy(userConfig, configInfo)
	if copyErr != nil {
		log.Error(strs.ObjectToString(copyErr))
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return configInfo, nil
}

func UpdateUserConfig(orgId, userId int64, input vo.UpdateUserConfigReq) (*vo.UpdateUserConfigResp, errs.SystemErrorInfo) {
	userConfigBo := &bo.UserConfigBo{}
	copyErr := copyer.Copy(input, userConfigBo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err1 := domain.UpdateUserConfig(orgId, userId, *userConfigBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err1)
	}

	//Nico: 暂时先这样做双写，之后优化
	err1 = domain.DeleteUserConfigInfo(orgId, userId)
	if err1 != nil {
		log.Error(err1)
	}

	return &vo.UpdateUserConfigResp{
		ID: input.ID,
	}, nil
}

func UpdateUserPcConfig(orgId, userId int64, input vo.UpdateUserPcConfigReq) (*vo.UpdateUserConfigResp, errs.SystemErrorInfo) {
	userConfigBo := &bo.UserConfigBo{}

	if util.FieldInUpdate(input.UpdateFields, "pcNoticeOpenStatus") {
		userConfigBo.PcNoticeOpenStatus = *input.PcNoticeOpenStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcIssueRemindMessageStatus") {
		userConfigBo.PcIssueRemindMessageStatus = *input.PcIssueRemindMessageStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcOrgMessageStatus") {
		userConfigBo.PcOrgMessageStatus = *input.PcOrgMessageStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcProjectMessageStatus") {
		userConfigBo.PcProjectMessageStatus = *input.PcProjectMessageStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcCommentAtMessageStatus") {
		userConfigBo.PcCommentAtMessageStatus = *input.PcCommentAtMessageStatus
	}

	err1 := domain.UpdateUserPcConfig(orgId, userId, *userConfigBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err1)
	}

	//Nico: 暂时先这样做双写，之后优化
	err1 = domain.DeleteUserConfigInfo(orgId, userId)
	if err1 != nil {
		log.Error(err1)
	}

	return &vo.UpdateUserConfigResp{
		ID: 0,
	}, nil
}

func UpdateUserDefaultProjectIdConfig(orgId, userId int64, input vo.UpdateUserDefaultProjectConfigReq) (*vo.UpdateUserConfigResp, errs.SystemErrorInfo) {
	userConfigBo, err := domain.GetUserConfigInfo(orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defaultProjectId := input.DefaultProjectID

	cacheProjectInfoResp := projectfacade.GetCacheProjectInfo(projectvo.GetCacheProjectInfoReqVo{
		OrgId:     orgId,
		ProjectId: defaultProjectId,
	})
	if cacheProjectInfoResp.Failure() {
		log.Error(cacheProjectInfoResp.Message)
		return nil, cacheProjectInfoResp.Error()
	}

	err1 := domain.UpdateUserDefaultProjectIdConfig(orgId, userId, *userConfigBo, defaultProjectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err1)
	}

	err1 = domain.DeleteUserConfigInfo(orgId, userId)
	if err1 != nil {
		log.Error(err1)
	}

	return &vo.UpdateUserConfigResp{
		ID: userConfigBo.ID,
	}, nil
}

func UpdateUserInfo(orgId, userId int64, input vo.UpdateUserInfoReq) (*vo.Void, errs.SystemErrorInfo) {

	if input.Name == nil && input.Sex == nil && input.Avatar == nil &&
		input.Birthday == nil && input.RemindBindPhone == nil {
		// 不再提醒弹窗
		return domain.SetRemindPopUp(orgId, userId, input.UpdateFields)
	}

	upd := &mysql.Upd{}
	//头像
	assemblyAvatar(input, upd)
	//姓名
	nameErr := assemblyName(input, upd)

	if nameErr != nil {
		log.Error(nameErr)
		return nil, nameErr
	}

	//出生日期
	assemblyBirthday(input, upd)
	//性别
	sexErr := assemblySex(input, upd)

	if sexErr != nil {
		log.Error(sexErr)
		return nil, sexErr
	}

	if NeedUpdate(input.UpdateFields, "remindBindPhone") {

		if input.RemindBindPhone != nil {

			if ok, _ := slice.Contain([]int{consts.AppIsRemind, consts.AppIsNotRemind}, *input.RemindBindPhone); ok {
				(*upd)[consts.TcRemindBindPhone] = *input.RemindBindPhone
			}
		}
	}

	if len(*upd) > 0 {
		err := domain.UpdateUserInfo(userId, *upd)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	err := domain.ClearBaseUserInfo(orgId, userId)
	if err != nil {
		log.Error(err)
	}

	return &vo.Void{
		ID: userId,
	}, nil
}

func assemblySex(input vo.UpdateUserInfoReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "sex") {

		if input.Sex != nil {

			if *input.Sex != consts.Male && *input.Sex != consts.Female {
				return errs.BuildSystemErrorInfo(errs.UserSexFail)
			}
			(*upd)[consts.TcSex] = *input.Sex
		}
	}
	return nil
}

func assemblyBirthday(input vo.UpdateUserInfoReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "birthday") {

		if input.Birthday != nil {
			birthday := time.Time(*input.Birthday)
			(*upd)[consts.TcBirthday] = birthday
		}
	}
}

// 组装个人头像信息
func assemblyAvatar(input vo.UpdateUserInfoReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "avatar") {

		if input.Avatar != nil {
			(*upd)[consts.TcAvatar] = *input.Avatar
		}
	}
}

// 组装名字
func assemblyName(input vo.UpdateUserInfoReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "name") {

		if input.Name != nil {

			name := strings.Trim(*input.Name, " ")
			//nameLen := str.CountStrByGBK(name)
			//
			//if nameLen == 0 || nameLen > 20 {
			//	log.Error("姓名长度错误")
			//	return errs.BuildSystemErrorInfo(errs.UserNameLenError)
			//}
			isNameRight := format.VerifyUserNameFormat(name)
			if !isNameRight {
				return errs.BuildSystemErrorInfo(errs.UserNameLenError)
			}
			(*upd)[consts.TcName] = name
			(*upd)[consts.TcNamePinyin] = pinyin.ConvertToPinyin(name)
		}
	}
	return nil
}

func GetUserInfoByUserIds(input orgvo.GetUserInfoByUserIdsReqVo) (*[]orgvo.GetUserInfoByUserIdsRespVo, errs.SystemErrorInfo) {
	bos, err := domain.GetBaseUserInfoBatch(input.OrgId, input.UserIds)
	if err != nil {
		return nil, err
	}

	vos := &[]orgvo.GetUserInfoByUserIdsRespVo{}

	copyError := copyer.Copy(bos, vos)
	if copyError != nil {
		log.Error(copyError)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyError)
	}
	return vos, nil
}

func VerifyOrg(orgId int64, userId int64) bool {
	return domain.VerifyOrg(orgId, userId)
}

func VerifyOrgUsers(orgId int64, userIds []int64) bool {
	return domain.VerifyOrgUsers(orgId, userIds)
}

func VerifyOrgUsersReturnValid(orgId int64, userIds []int64) []int64 {
	return domain.VerifyOrgUsersReturnValid(orgId, userIds)
}

func GetUserInfo(orgId int64, userId int64, sourceChannel string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	res, _, err := domain.GetUserInfo(orgId, userId, sourceChannel)
	return res, err
}

func GetOutUserInfoListBySourceChannel(sourceChannel string, page int, size int) ([]bo.UserOutInfoBo, errs.SystemErrorInfo) {
	return domain.GetOutUserInfoListBySourceChannel(sourceChannel, page, size)
}

func GetOutUserInfoListByUserIds(idList []int64) ([]bo.UserOutInfoBo, errs.SystemErrorInfo) {
	return domain.GetOutUserInfoListByUserIds(idList)
}

func GetUserInfoListByOrg(orgId int64) ([]bo.SimpleUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetUserInfoListByOrg(orgId)
}

func JudgeUserIsAdmin(outOrgId string, outUserId string, sourceChannel string) bool {
	return domain.JudgeUserIsAdmin(outOrgId, outUserId, sourceChannel)
}

// 异步的方式，推送任务到mq，进行同步成员、部门信息
func PushSyncMemberDept(orgId, currentUserId int64, syncUserInfoFromFeiShu vo.SyncUserInfoFromFeiShuReq) errs.SystemErrorInfo {
	// 用户配置了权限，才能调用后续的逻辑
	authErr := AuthOrgRole(orgId, currentUserId, consts.RoleOperationPathOrgUser, consts.RoleOperationModifyDepartment)
	if authErr != nil {
		return authErr
	}
	param := orgvo.SyncUserInfoFromFeiShuReqVo{
		CurrentUserId:             currentUserId,
		OrgId:                     orgId,
		SyncUserInfoFromFeiShuReq: syncUserInfoFromFeiShu,
	}
	return domain.PushSyncMemberDept(param)
}

// 从飞书方更新成员信息
// 1.获取当前组织下的所有成员
// 2.通过成员列表的信息请求获取飞书方的成员信息，更新我方成员信息
// 3.删除旧部门、同步新部门信息
func SyncUserInfoFromFeiShu(orgId, currentUserId int64, syncUserInfoFromFeiShu vo.SyncUserInfoFromFeiShuReq) errs.SystemErrorInfo {
	baseOrgOutInfo, err := domain.GetBaseOrgOutInfo(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	corpId := baseOrgOutInfo.OutOrgId
	// 获取被授权的部门列表
	outDeptList, err := feishu.GetScopeDeps(corpId)
	// 获取组织下所有成员列表
	userInfoBoList, err := domain.GetUserInfoListByOrg(orgId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.UserNotFoundError, err)
	}
	var userIdList []int64
	for _, val := range userInfoBoList {
		userIdList = append(userIdList, val.Id)
	}
	// 查询成员的 openid 等外部信息，
	outInfoBos, err := domain.GetOutUserInfoListByUserIds(userIdList)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.UserNotFoundError, err)
	}
	// 获取 openid list
	var mapOpenIdToUserId = map[string]int64{}
	for _, outInfoBo := range outInfoBos {
		mapOpenIdToUserId[outInfoBo.OutOrgUserId] = outInfoBo.UserId
	}
	// 这里通过获取作用域内的所有用户
	// 注意：outOpenIdList 中存在着大量的 outUserId，获取其信息时，需要分批获取！
	outOpenIdList, err := feishu.GetScopeOpenIds(corpId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.FeiShuUserNotInAppUseScopeOfAuthority, err)
	}
	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := updateUserFromFsWithPagination(orgId, currentUserId, corpId, outOpenIdList, syncUserInfoFromFeiShu, mapOpenIdToUserId, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		// 是否要同步部门数据
		if syncUserInfoFromFeiShu.NeedSyncDepartment {
			oriErr := syncDepartments(orgId, outDeptList, corpId, tx)
			if oriErr != nil {
				return errs.BuildSystemErrorInfo(errs.SyncDepartmentError, err)
			}
		}
		return nil
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return nil
}

// 传入一批 openId，分批次对其更新处理
func updateUserFromFsWithPagination(orgId, currentUserId int64, corpId string, batchOpenIds []string, syncUserInfoFromFeiShu vo.SyncUserInfoFromFeiShuReq, mapOpenIdToUserId map[string]int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	// 调用飞书 sdk 获取飞书的用户信息
	tenant, err := feishu.GetTenant(corpId)
	if err != nil {
		log.Error(err)
		return err
	}
	batch := 50
	offset := 0
	surplusSize := len(batchOpenIds)
	if surplusSize > 0 {
		for {
			limit := offset + batch
			if surplusSize < limit {
				limit = surplusSize
			}
			tmpOpenIds := batchOpenIds[offset:limit]
			userBatchResp, oriErr := tenant.GetUserBatchGetV2(nil, tmpOpenIds)
			if oriErr != nil {
				log.Error(oriErr)
				return errs.BuildSystemErrorInfo(errs.FeiShuClientTenantError, oriErr)
			}
			// 更新员工信息：姓名、头像
			oriErr = updateUserInfoFromFeiShu(orgId, currentUserId, corpId, syncUserInfoFromFeiShu, mapOpenIdToUserId, userBatchResp.Data.Users, tx)
			if oriErr != nil {
				log.Error(oriErr)
				return errs.BuildSystemErrorInfo(errs.SyncUserUpdateUserFail, err)
			}
			if surplusSize <= limit {
				break
			}
			offset += batch
		}
	}
	return nil
}

// fsUsers 从飞书获取的用户信息，更新用户数据到极星
// 每更新一个 sleep 10ms，防止出现长时间占用数据库的锁
func updateUserInfoFromFeiShu(orgId, currentUserId int64, outOrgId string, syncUserInfoFromFeiShu vo.SyncUserInfoFromFeiShuReq, mapOpenIdToUserId map[string]int64, fsUsers []sdkVo.UserDetailInfoV2, tx sqlbuilder.Tx) error {
	// 如果没有要更新的用户 out 信息，则直接返回
	if len(fsUsers) < 1 {
		return nil
	}
	var outUserInfoBosKeyByOpenId = map[string]sdkVo.UserDetailInfoV2{}
	var tmpUpdateReq vo.UpdateUserInfoReq
	var tmpUpd mysql.Upd
	var err error
	var userIdList = []int64{}
	var needInsertOpenIdMap = map[string]sdkVo.UserDetailInfoV2{}
	var needInsertOpenUsersDeptMap = map[string]string{}
	var needInsertUsersDeptOutIdList = []string{}
	for _, oneFsUser := range fsUsers {
		oneFsUser.Name = strings.Trim(oneFsUser.Name, " ")
		if len(oneFsUser.Name) < 1 {
			oneFsUser.Name = "未命名"
		}
		outUserInfoBosKeyByOpenId[oneFsUser.OpenId] = oneFsUser
		tmpUpdateReq = vo.UpdateUserInfoReq{
			Name:            nil,
			Sex:             nil,
			Avatar:          nil,
			Birthday:        nil,
			RemindBindPhone: nil,
			UpdateFields:    []string{},
		}
		tmpUpd = mysql.Upd{}
		// 根据传入的 option 参数决定是否要更新对应的数据 todo
		if syncUserInfoFromFeiShu.NeedSyncName {
			tmpUpdateReq.UpdateFields = append(tmpUpdateReq.UpdateFields, consts.TcName)
			tmpUpdateReq.Name = &oneFsUser.Name
			tmpUpd[consts.TcName] = tmpUpdateReq.Name
		}
		if syncUserInfoFromFeiShu.NeedSyncAvatar && len(oneFsUser.Avatar.AvatarOrigin) > 0 {
			tmpUpdateReq.UpdateFields = append(tmpUpdateReq.UpdateFields, consts.TcAvatar)
			tmpUpdateReq.Avatar = &oneFsUser.Avatar.AvatarOrigin
			tmpUpd[consts.TcAvatar] = tmpUpdateReq.Avatar
		}
		if tmpUserId, ok := mapOpenIdToUserId[oneFsUser.OpenId]; ok {
			// 更新 ppm_org_user 和 ppm_org_user_out_info
			_, err = UpdateUserInfo(orgId, tmpUserId, tmpUpdateReq)
			if err != nil {
				return err
			}
			err = domain.UpdateUserOutInfo(orgId, oneFsUser.OpenId, tmpUpd)
			if err != nil {
				return err
			}
			userIdList = append(userIdList, tmpUserId)
		} else {
			// 在极星应用中，找不到该用户，此时需要将这个用户同步过来
			// 先暂时把该用户id 存储一下，后续执行新增新的用户
			needInsertOpenIdMap[oneFsUser.OpenId] = oneFsUser
			for _, tmpOutDeptId := range oneFsUser.Departments {
				needInsertOpenUsersDeptMap[oneFsUser.OpenId] = tmpOutDeptId
				needInsertUsersDeptOutIdList = append(needInsertUsersDeptOutIdList, tmpOutDeptId)
				break
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
	if len(needInsertOpenIdMap) > 0 {
		for openId, _ := range needInsertOpenIdMap {
			// 通过 openId outOrgId 等信息，新增用户
			_, err := domain.InitPlatformUser(sdk_const.SourceChannelFeishu, orgId, outOrgId, openId, tx, "")
			if err != nil {
				return err
			}
		}
	}
	if len(userIdList) < 1 {
		return nil
	}
	//最后将用户信息缓存清掉
	clearErr := domain.ClearBaseUserInfoBatch(orgId, userIdList)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return err
}

// delete(soft) departments by out departmentId
// （软）删除旧部门相关数据。涉及到 ppm_org_department、ppm_org_department_out_info、ppm_org_user_department 表
func deleteDepartments(orgId int64, outDepartmentIds []string) error {
	// 如果没有要删除的部门，则直接返回
	if orgId < 1 {
		return nil
	}
	// 删除部门的 out 信息
	err := domain.DeleteDepartmentOutInfoByOrgId(orgId)
	if err != nil {
		return err
	}
	// 删除旧部门
	err = domain.DeleteDepartmentByOrgId(orgId)
	if err != nil {
		return err
	}
	// 删除用户和部门的关联关系
	err = domain.DeleteUserDepartmentByOrgId(orgId)
	if err != nil {
		return err
	}
	return err
}

// sync department list
// 1.删除原先的所有部门
// 2.同步飞书方的部门信息
func syncDepartments(orgId int64, outDeptList []sdkVo.DepartmentRestInfoVo, tenantKey string, tx sqlbuilder.Tx) error {
	depSize := len(outDeptList)
	var (
		outDeptIdList []string
		fsDepIdMap    = map[string]int64{}
	)
	// 生成新的部门id，和新的 out 部门信息
	depIds, err1 := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartment, depSize)
	if err1 != nil {
		log.Error(err1)
		return errors.New(err1.Message())
	}
	depOutIds, err1 := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartmentOutInfo, depSize)
	if err1 != nil {
		log.Error(err1)
		return errors.New(err1.Message())
	}
	// 获取授权的部门列表，即 outDeptList
	for k, oneDept := range outDeptList {
		outDeptIdList = append(outDeptIdList, oneDept.Id)
		// 构建 outDepartmentId => departmentId 的 map
		fsDepIdMap[oneDept.Id] = depIds.Ids[k].Id
	}
	// 删除原先旧的部门信息
	err := deleteDepartments(orgId, outDeptIdList)
	if err != nil {
		return err
	}
	// 构建新增部门的数组，以及生成 outDeptId => 我方部门id 的 map2
	newDepartmentList, newOutDepartmentList, err := AssemblyNewDeptDataList(orgId, outDeptList, fsDepIdMap, depIds, depOutIds, depSize)
	err = mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, newDepartmentList)
	if err != nil {
		log.Error(err)
		return err
	}
	err = mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, newOutDepartmentList)
	if err != nil {
		log.Error(err)
		return err
	}
	// 通过多个部门ids，获取其下的所有员工
	surplusOpenIds, err1 := feishu.GetScopeOpenIds(tenantKey)
	if err != nil {
		log.Error(err)
		return errors.New(err1.Message())
	}
	surplusSize := len(surplusOpenIds)
	// 如果授权管理下，没有员工，则直接返回
	if surplusSize < 1 {
		return errors.New(errs.SyncUserHasNoUserUnderPermission.Message())
	}
	tenant, err1 := feishu.GetTenant(tenantKey)
	if err1 != nil {
		log.Error(err1)
		return errors.New(err1.Message())
	}
	// 构建更新用户与部门的关联关系数据
	userDepList, err := AssemblyUserDeptDataList(orgId, tenant, fsDepIdMap, surplusOpenIds)
	// 新增用户和部门的关联关系
	if len(userDepList) > 0 {
		userDepPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(userDepList))
		if idErr != nil {
			log.Error(idErr)
			return errors.New(idErr.Message())
		}
		for i, _ := range userDepList {
			userDepList[i].Id = userDepPoIds.Ids[i].Id
		}
		err = mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, slice.ToSlice(userDepList))
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

// 组装新增部门的数据。部门列表，部门 out 列表
func AssemblyNewDeptDataList(orgId int64, outDeptList []sdkVo.DepartmentRestInfoVo, fsDepIdMap map[string]int64, depIds *bo.IdCodes, depOutIds *bo.IdCodes, depSize int) ([]interface{}, []interface{}, error) {
	var (
		newDepartmentList    = make([]interface{}, depSize)
		newOutDepartmentList = make([]interface{}, depSize)
	)
	rootId := fsDepIdMap["0"]
	for k, v := range outDeptList {
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
		newDepartmentList[k] = &po.PpmOrgDepartment{
			Id:            depId,
			OrgId:         orgId,
			Name:          v.Name,
			ParentId:      parentDepId,
			SourceChannel: sdk_const.SourceChannelFeishu,
		}
		newOutDepartmentList[k] = po.PpmOrgDepartmentOutInfo{
			Id:                       depOutId,
			OrgId:                    orgId,
			DepartmentId:             depId,
			SourceChannel:            sdk_const.SourceChannelFeishu,
			OutOrgDepartmentId:       v.Id,
			Name:                     v.Name,
			OutOrgDepartmentParentId: v.ParentId,
		}
	}
	return newDepartmentList, newOutDepartmentList, nil
}

// 构建更新用户与部门的关联关系数据，用于后续的入库
func AssemblyUserDeptDataList(orgId int64, tenant *sdk.Tenant, fsDepIdMap map[string]int64, surplusOpenIds []string) ([]po.PpmOrgUserDepartment, error) {
	surplusSize := len(surplusOpenIds)
	batch := 30
	offset := 0
	// 用于去重，防止下方的循环处理中，对一个用户重复组装数据
	userIdCache := map[string]bool{}
	userDepList := make([]po.PpmOrgUserDepartment, 0)
	rootDepId := fsDepIdMap["0"]
	for {
		limit := offset + batch
		if surplusSize < limit {
			limit = surplusSize
		}
		openIds := surplusOpenIds[offset:limit]
		userBatchResp, err := tenant.GetUserBatchGetV2(nil, openIds)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if userBatchResp.Code != 0 {
			log.Error(userBatchResp.Msg)
			return nil, err
		}
		userPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUser, len(openIds))
		if idErr != nil {
			log.Error(idErr)
			return nil, errors.New(idErr.Message())
		}
		for i, fsUser := range userBatchResp.Data.Users {
			userNativeId := userPoIds.Ids[i].Id
			userDeps := fsUser.Departments
			if userDeps != nil && len(userDeps) > 0 {
				hasDep := false
				for _, userDep := range userDeps {
					if depId, ok := fsDepIdMap[userDep]; ok {
						//这些用户归入部门
						userDepList = append(userDepList, AssemblyFeiShuUserDepRelationInfo(orgId, userNativeId, depId))
						hasDep = true
					}
				}
				// 如果没有所属部门，则归类为顶级部门（公司）下
				if !hasDep {
					userDepList = append(userDepList, AssemblyFeiShuUserDepRelationInfo(orgId, userNativeId, rootDepId))
				}
			}
			userIdCache[fsUser.OpenId] = true
		}
		if surplusSize <= limit {
			break
		}
		offset += batch
	}

	return userDepList, nil
}

func AssemblyFeiShuUserOutInfo(orgId int64, userId int64, fsUserDetailInfo sdkVo.UserDetailInfoV2) po.PpmOrgUserOutInfo {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOutInfo := &po.PpmOrgUserOutInfo{}
	userOutInfo.UserId = userId
	userOutInfo.OrgId = orgId
	userOutInfo.OutOrgUserId = fsUserDetailInfo.OpenId
	userOutInfo.OutUserId = fsUserDetailInfo.OpenId
	userOutInfo.IsDelete = consts.AppIsNoDelete
	userOutInfo.Status = consts.AppStatusEnable
	userOutInfo.SourceChannel = sdk_const.SourceChannelFeishu
	userOutInfo.Name = fsUserDetailInfo.Name
	userOutInfo.Avatar = fsUserDetailInfo.Avatar.AvatarOrigin
	userOutInfo.JobNumber = fsUserDetailInfo.EmployeeNo

	return *userOutInfo
}

func AssemblyFeiShuUserOrgRelationInfo(orgId int64, userId int64) po.PpmOrgUserOrganization {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOrgRelationInfo := &po.PpmOrgUserOrganization{}
	userOrgRelationInfo.UserId = userId
	userOrgRelationInfo.OrgId = orgId
	userOrgRelationInfo.Status = consts.AppStatusEnable
	userOrgRelationInfo.UseStatus = consts.AppStatusDisabled
	userOrgRelationInfo.CheckStatus = consts.AppCheckStatusSuccess

	return *userOrgRelationInfo
}

func AssemblyFeiShuUserDepRelationInfo(orgId int64, userId int64, depId int64) po.PpmOrgUserDepartment {
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userDepRelationInfo := &po.PpmOrgUserDepartment{}
	userDepRelationInfo.UserId = userId
	userDepRelationInfo.OrgId = orgId
	userDepRelationInfo.DepartmentId = depId

	return *userDepRelationInfo
}

func GetUserInfoByFeishuTenantKey(tenantKey string, openId string) (orgvo.GetUserInfoByFeishuTenantKeyData, errs.SystemErrorInfo) {
	result := orgvo.GetUserInfoByFeishuTenantKeyData{
		OrgId:  0,
		UserId: 0,
	}
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOutOrgId: tenantKey,
	}, orgOutInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, errs.OrgNotExist
		} else {
			log.Error(err)
			return result, errs.MysqlOperateError
		}
	}
	result.OrgId = orgOutInfo.OrgId

	outUserInfo := &po.PpmOrgUserOutInfo{}
	userErr := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgOutInfo.OrgId,
		consts.TcOutUserId: openId,
	}, outUserInfo)
	if userErr != nil {
		if err == db.ErrNoMoreRows {
			return result, errs.UserNotExist
		} else {
			log.Error(err)
			return result, errs.MysqlOperateError
		}
	}
	result.UserId = outUserInfo.UserId

	return result, nil
}

func GetOrgUserIdsByEmIds(orgId int64, sourceChannel string, empIds []string) (map[string]int64, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgUserOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcOutUserId:     db.In(empIds),
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := map[string]int64{}
	for _, info := range *pos {
		res[info.OutUserId] = info.UserId
	}

	return res, nil
}

// 通过一批用户id，返回这批id中离职的用户id。fs应用下的离职的用户，由 ppm_org_user_organization 表中的 is_delete 标识。
func FilterResignedUserIds(orgId int64, userIds []int64) ([]int64, errs.SystemErrorInfo) {
	validUserIds := []int64{}
	if len(userIds) < 1 {
		return validUserIds, nil
	}
	pos := &[]po.PpmOrgUserOrganization{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, val := range *pos {
		validUserIds = append(validUserIds, val.UserId)
	}
	return validUserIds, nil
}

func GetOrgUserIds(orgId int64) ([]int64, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgUserOrganization{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	result := []int64{}
	for _, organization := range *pos {
		result = append(result, organization.UserId)
	}
	return result, nil
}

func GetPayRemind(orgId, userId int64) (*vo.GetPayRemindResp, errs.SystemErrorInfo) {
	remindMsg := ""
	//判断当前用户是不是超管
	//adminFlagBo, err := rolefacade.GetUserAdminFlagRelaxed(orgId, userId)
	//if err != nil{
	//	log.Error(err)
	//	return nil, err
	//}
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	if manageAuthInfoResp.Failure() {
		log.Error(manageAuthInfoResp.Message)
		return nil, manageAuthInfoResp.Error()
	}
	adminFlagBo := manageAuthInfoResp.Data
	//判断是否需要提示付费到期信息
	info, infoErr := domain.GetOrgConfig(orgId)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, infoErr
	}
	payName := businees.GetOrgPayVersionName(info.PayLevel)
	if info.PayStartTime.AddDate(0, 0, 16).After(info.PayEndTime) {
		//时间小于16天，判断为试用版
		// payName = "试用版"
	}
	now := time.Now()
	if businees.CheckIsPaidVer(info.PayLevel) && now.AddDate(0, 0, 3).After(info.PayEndTime) {
		//付费期到期前三天内
		if adminFlagBo.IsSysAdmin {
			needRemind, err := domain.NeedRemindPayExpire(orgId, userId, info.PayEndTime)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			if needRemind {
				expireTime := info.PayEndTime.Unix() - now.Unix()
				day := expireTime / 86400
				if expireTime%86400 > 0 {
					day += 1
				}
				remindMsg = fmt.Sprintf("%s将于%d天后到期，如需继续使用%s可前往管理后台https://feishu.cn/admin/index购买", payName, day, payName)
			}
		}
	} else if info.PayLevel == consts.PayLevelStandard && now.AddDate(0, 0, -1).Before(info.PayEndTime) && now.After(info.PayEndTime) && info.PayStartTime.Before(info.PayEndTime) {
		//付费之后一天内
		needRemind, err := domain.NeedRemindPayOverdue(orgId, userId, info.PayEndTime)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if needRemind {
			if adminFlagBo.IsSysAdmin {
				remindMsg = fmt.Sprintf("%s已到期，现已为您切换至免费版，如需再次购买%s请前往飞书管理后台https://feishu.cn/admin/index，选择“应用付费—极星协作”后购买需要的人数及使用时长",
					payName, payName)
			}
			//普通用户不提示
			//} else {
			//	remindMsg = fmt.Sprintf("%s已到期，现已为您切换至免费版，如需继续使用标准版可联系管理员前往管理后台https://feishu.cn/admin/index购买", payName)
			//}
		}
	}

	return &vo.GetPayRemindResp{RemindPayExpireMsg: remindMsg}, nil
}

func GetOrgUsersInfoByEmIds(orgId int64, sourceChannel string, empIds []string) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgUserOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcOutUserId:     db.In(empIds),
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	userIds := []int64{}
	for _, info := range *pos {
		userIds = append(userIds, info.UserId)
	}

	return domain.GetBaseUserInfoBatch(orgId, userIds)
}

func GetUserOutInfoByOpenID(orgId int64, sourceChannel string, openId string) (*po.PpmOrgUserOutInfo, errs.SystemErrorInfo) {
	//pos := &po.PpmOrgUserOutInfo{}
	//err := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
	//	consts.TcIsDelete:      consts.AppIsNoDelete,
	//	consts.TcSourceChannel: sourceChannel,
	//	consts.TcOutUserId:     openId,
	//}, pos)
	//默认去取最新的一条，历史数据或者前端缓存可能导致初始化了out_info的org_id为0的用户，这样能够有效避免
	pos := []po.PpmOrgUserOutInfo{}
	count, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: sourceChannel,
		consts.TcOutUserId:     openId,
		consts.TcOrgId:         orgId,
	}, nil, 0, 1, "id desc", &pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	if count == 0 {
		return nil, errs.UserOutInfoNotExist
	}

	return &pos[0], nil
}

func GetUserOutInfoByPlatformAndOrgId(sourceChannel string, userId, orgId int64) ([]po.PpmOrgUserOutInfo, errs.SystemErrorInfo) {
	pos := make([]po.PpmOrgUserOutInfo, 0)
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        userId,
		consts.TcOrgId:         orgId,
	}, &pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	return pos, nil
}

func GetOrgUserAndDeptCount(orgId int64) (*orgvo.GetOrgUserAndDeptCountRespData, errs.SystemErrorInfo) {
	userCount, mysqlErr := mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	})
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}
	deptCount, mysqlErr := mysql.SelectCountByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}
	return &orgvo.GetOrgUserAndDeptCountRespData{
		UserCount: userCount,
		DeptCount: deptCount,
	}, nil
}

// 设定用户浏览了“新手指引”状态
func SetVisitUserGuideStatus(input orgvo.SetVisitUserGuideStatusReq) (bool, errs.SystemErrorInfo) {
	// 模板预览的组织不展示“新手指南”
	if input.OrgId == consts.PreviewTplOrgId {
		return true, nil
	}
	userConfigBo, err := domain.GetUserConfigInfo(input.OrgId, input.UserId)
	if err != nil {
		if strings.Contains(err.Error(), db.ErrNoMoreRows.Error()) {
			// 没有则新增 user config 记录
			userConfigBo, err = domain.InsertUserConfig(input.OrgId, input.UserId)
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else {
			log.Error(err)
			return false, err
		}
	}
	userConfigExt := bo.UserConfigExt{}
	if userConfigBo.Ext != "" {
		errJson := json.FromJson(userConfigBo.Ext, &userConfigExt)
		if errJson != nil {
			log.Errorf("[SetVisitUserGuideStatus]err:%v", errJson)
			return false, errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
		}
	}
	// 设置新手指引flag时设置版本信息弹窗不展示
	version := &bo.VersionConfig{VersionInfoVisible: false}
	userConfigExt.Version = version
	switch input.Flag {
	case "newUserGuideStatus": // 用户指引
		userConfigExt.VisitedNewUserGuide = true
	default:
		return false, errs.ParamError
	}
	newExtJson := json.ToJsonIgnoreError(userConfigExt)
	if err = domain.UpdateUserConfig(input.OrgId, input.UserId, bo.UserConfigBo{
		ID:  userConfigBo.ID,
		Ext: newExtJson,
	}); err != nil {
		log.Error(err)
		return false, err
	}
	// 更新之后需要清理一下 user config 缓存
	cacheKey, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  input.OrgId,
		consts.CacheKeyUserIdConstName: input.UserId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	if _, err := cache.Del(cacheKey); err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return true, nil
}

func SetVersionVisible(req orgvo.SetVersionReq) (bool, errs.SystemErrorInfo) {
	err := domain.SetVersionConfig(req.OrgId, req.UserId, req.Input.VersionInfoVisible)
	if err != nil {
		log.Errorf("[SetVersionConfig]错误, err: %v", err)
		return false, err
	}
	return true, nil
}

func SetUserActivity(req orgvo.SetUserActivityReq) (bool, errs.SystemErrorInfo) {
	err := domain.SetActivity11Info(req.OrgId, req.UserId, req.Input.ActivityFlag)
	if err != nil {
		log.Errorf("[SetUserActivity] err:%v", err)
		return false, err
	}
	return true, nil
}

func GetVersion(req orgvo.GetVersionReq) (*orgvo.VersionResp, errs.SystemErrorInfo) {
	cardInfo, err := domain.GetCardConfig(req.OrgId, req.UserId)
	if err != nil {
		return nil, err
	}
	return cardInfo, nil
}

func SetUserViewLocation(req orgvo.SaveViewLocationReqVo) (bool, errs.SystemErrorInfo) {
	err := domain.SaveUserViewLocation(req)
	if err != nil {
		log.Errorf("[SetUserViewLocation] err:%v", err)
		return false, err
	}
	return true, nil
}

func GetUserViewLocation(req orgvo.GetViewLocationReq) ([]*orgvo.UserLastViewLocationData, errs.SystemErrorInfo) {
	userViewLocations, err := domain.GetUserViewLocation(req.OrgId, req.UserId)
	if err != nil {
		log.Errorf("[GetUserViewLocation] err:%v", err)
		return nil, err
	}
	data := []*orgvo.UserLastViewLocationData{}
	errCopy := copyer.Copy(userViewLocations, &data)
	if errCopy != nil {
		log.Errorf("[GetUserViewLocation] copy err:%v", errCopy)
		return nil, errs.ObjectCopyError
	}
	return data, nil
}

func GetViewLocationList(req orgvo.GetViewLocationReq) ([]*orgvo.UserLastViewLocationData, errs.SystemErrorInfo) {
	viewLocationList, err := domain.GetViewLocationList(req.OrgId, req.UserId)
	if err != nil {
		log.Errorf("[GetViewLocationList]err:%v", err)
		return nil, err
	}
	return viewLocationList, nil
}

// 删除应用的时候，把用户浏览应用的位置信息去除 DeleteAppUserLocation
func DeleteAppUserLocation(orgId, userId, appId int64) errs.SystemErrorInfo {
	return domain.DeleteUserLocationWithAppId(orgId, userId, appId)
}

// createOrgUserByAccount 用户名添加成员
func createOrgUserByAccount(orgId, operatorId int64, req orgvo.CreateOrgMemberReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	// 验证账号名
	if !format.VerifyUserNameFormat(req.AccountName) {
		return nil, errs.AccountNameLenError
	}
	// 检查账号是否已存在
	if domain.CheckUserAccountByLoginName(orgId, req.AccountName) {
		return nil, errs.AccountAllReadyExist
	}

	userInfoBo, err := domain.UserRegister(bo.UserSMSRegisterInfo{
		OrgId:          orgId,
		SourceChannel:  consts.AppSourceChannelWeb,
		SourcePlatform: consts.AppSourceChannelWeb,
		AccountName:    req.AccountName,
		Name:           req.Name,
		Password:       req.Password,
		Status:         req.Status,
	})
	if err != nil {
		log.Errorf("[CreateOrgUserByAccount] UserRegister err:%v, orgId:%v, accountName:%v", err, orgId, req.AccountName)
		return nil, err
	}
	err = domain.AddOrgMember(orgId, userInfoBo.ID, operatorId, false, false)
	if err != nil {
		log.Errorf("[CreateOrgUserByAccount] AddOrgMember err:%v, orgId:%v, accountName:%v", err, orgId, req.AccountName)
		return nil, err
	}

	return userInfoBo, nil
}

func createOrgUserByPhoneNumber(orgId, operator int64, req orgvo.CreateOrgMemberReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	if !format.VerifyUserNameFormat(req.Name) {
		return nil, errs.UserNameLenError
	}
	phoneInvalid, errNum := phonenumbers.Parse(req.PhoneRegion+req.PhoneNumber, "")
	if errNum != nil || !phonenumbers.IsValidNumber(phoneInvalid) {
		return nil, errs.MobileInvalidError
	}
	phoneNumber := fmt.Sprintf("%s-%s", req.PhoneRegion, req.PhoneNumber)
	needRegisterLoginNames, needCreateOrgUserRelationUserIdsMap, needResetCheckStatusUserMap, errSys := domain.DetectUserInfoInUser(orgId, []string{phoneNumber})
	if errSys != nil {
		log.Errorf("[CreateOrgUser] DetectUserInfoInUser err:%v, orgId:%v", errSys, orgId)
		return nil, errSys
	}

	if len(needRegisterLoginNames) == 0 && len(needCreateOrgUserRelationUserIdsMap) == 0 && len(needResetCheckStatusUserMap) == 0 {
		return nil, errs.MobileSameError
	}

	inCheck := false
	inDisabled := false
	if req.Status == consts.AppStatusDisabled {
		inCheck = false
		inDisabled = true
	}

	var userBo *bo.UserInfoBo
	if len(needRegisterLoginNames) > 0 {
		// 注册
		log.Infof("用户%s未注册，开始注册....", req.PhoneNumber)
		userBo, errSys = domain.UserRegister(bo.UserSMSRegisterInfo{
			OrgId:          orgId,
			PhoneNumber:    phoneNumber,
			SourceChannel:  consts.AppSourceChannelWeb,
			SourcePlatform: consts.AppSourceChannelWeb,
			Name:           req.Name,
			MobileRegion:   req.PhoneRegion,
			Password:       req.Password,
			Status:         req.Status,
		})
		if errSys != nil {
			log.Errorf("[CreateOrgUser] UserRegister err:%v, orgId:%v, phoneNumber:%v", errSys, orgId, phoneNumber)
			return nil, errSys
		}
	}

	if len(needResetCheckStatusUserMap) > 0 {
		memberIds := make([]int64, 0, len(needResetCheckStatusUserMap))
		for userId := range needResetCheckStatusUserMap {
			memberIds = append(memberIds, userId)
			userBo = &bo.UserInfoBo{ID: userId}
			break
		}
		errSys = domain.ModifyOrgMemberCheckStatus(orgId, memberIds, consts.AppCheckStatusSuccess, operator, true)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] ModifyOrgMemberCheckStatus err:%v, orgId:%v, phoneNumber:%v", errSys, orgId, phoneNumber)
			return nil, errSys
		}
	}

	if len(needCreateOrgUserRelationUserIdsMap) > 0 {
		// user 已存在，只需要增加用户和组织的关联
		for userId := range needCreateOrgUserRelationUserIdsMap {
			userBo = &bo.UserInfoBo{ID: userId}
			break
		}
		//// 修改名字
		//errSys = domain.UpdateOrgMemberInfo(orgId, userBo.ID, mysql.Upd{consts.TcName: req.Name}, nil)
		//if errSys != nil {
		//	log.Errorf("[CreateOrgUser] UpdateOrgMemberInfo err:%v, orgId:%v, userId:%v", errSys, orgId, userBo.ID)
		//	return errSys
		//}
	}

	// needRegisterLoginNames 和 needCreateOrgUserRelationUserIdsMap 两种场景 需要绑定org_user表
	if len(needRegisterLoginNames) > 0 || len(needCreateOrgUserRelationUserIdsMap) > 0 {
		errSys = domain.AddOrgMember(orgId, userBo.ID, operator, inCheck, inDisabled)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] AddOrgMember err:%v, orgId:%v, userId:%v", errSys, orgId, userBo.ID)
			return nil, errSys
		}
	}

	if len(needCreateOrgUserRelationUserIdsMap) > 0 || len(needResetCheckStatusUserMap) > 0 {
		// 修改密码，替换该手机号在所有团队中的密码
		if req.Password != "" {
			errSys = RetrievePassword(orgvo.RetrievePasswordReqVo{Input: vo.RetrievePasswordReq{
				Username:    phoneNumber,
				NewPassword: req.Password,
			}})
			if errSys != nil {
				log.Errorf("[CreateOrgUser] RetrievePassword err:%v, orgId:%v, phoneNumber:%v", errSys, orgId, phoneNumber)
				return nil, errSys
			}
		}
		// 修改名字
		// 设置了姓名的，仅替换该账号所有本地团队的姓名
		//userOrgIdsMap, err := domain.GetLocalOrgUserIdMap([]int64{userBo.ID})
		//if err != nil {
		//	log.Errorf("[CreateOrgUser] GetLocalOrgUserIdMap err:%v, orgId:%v, phoneNumber:%v", errSys, orgId, phoneNumber)
		//	return nil, err
		//}
		//for userId, orgIds := range userOrgIdsMap {
		//	errSys = domain.UpdateLocalOrgUserNames(orgIds, userId, req.Name)
		//	if errSys != nil {
		//		log.Errorf("[CreateOrgUser] UpdateLocalOrgUserNames err:%v, orgId:%v, userId:%v", errSys, orgId, userId)
		//		return nil, errSys
		//	}
		//}
	}

	return userBo, nil
}

// 添加组织成员
func CreateOrgUser(orgId, operator int64, req orgvo.CreateOrgMemberReq) errs.SystemErrorInfo {
	// 校验角色 (其实是管理组角色)
	if len(req.RoleGroupIds) < 1 {
		return errs.ReqParamsValidateError
	}
	manageGroupList, errSys := domain.GetManageGroupList(orgId, req.RoleGroupIds)
	if errSys != nil {
		log.Errorf("[CreateOrgUser] GetManageGroupList err:%v, orgId:%v", errSys, orgId)
		return errSys
	}
	if len(manageGroupList) != len(req.RoleGroupIds) {
		return errs.ManageGroupNotExist
	}
	// 校验部门
	if len(req.DepartmentIds) > 0 {
		depBos, errSys := domain.GetDepartmentBoListByIds(orgId, req.DepartmentIds)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] GetDepartmentBoListByIds err:%v, orgId:%v", errSys, orgId)
			return errSys
		}
		if len(*depBos) != len(req.DepartmentIds) {
			return errs.DepartmentNotExist
		}
	}

	var userId int64
	if req.AccountType == consts.AccountLogin {
		userBo, errSys := createOrgUserByAccount(orgId, operator, req)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] createOrgUserByAccount err:%v, orgId:%v, accountName:%v", errSys, orgId, req.AccountName)
			return errSys
		}
		userId = userBo.ID
	} else {
		userBo, errSys := createOrgUserByPhoneNumber(orgId, operator, req)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] createOrgUserByPhoneNumber err:%v, orgId:%v, accountName:%v", errSys, orgId, req.AccountName)
			return errSys
		}
		userId = userBo.ID
	}

	if len(req.RoleGroupIds) > 0 {
		errSys = domain.BindUserRoleGroups(orgId, operator, req.RoleGroupIds, []int64{userId})
		if errSys != nil {
			log.Errorf("[CreateOrgUser] BindUserRoleGroups err:%v, orgId:%v, userId:%v", errSys, orgId, userId)
			return errSys
		}
	}

	if len(req.DepartmentIds) > 0 {
		errSys = domain.BindUserDepartments(orgId, userId, req.DepartmentIds)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] BindUserDepartments err:%v, orgId:%v, userId:%v", errSys, orgId, userId)
			return errSys
		}
	}
	return nil
}

// 修改成员信息
func UpdateOrgUser(orgId, operator int64, req orgvo.UpdateOrgMemberReq) errs.SystemErrorInfo {
	if req.UserId == 0 {
		return errs.ReqParamsValidateError
	}
	if req.Name != "" {
		if !format.VerifyUserNameFormat(req.Name) {
			return errs.UserNameLenError
		}
		req.Name = strings.TrimSpace(req.Name)
	}
	// 非超管不能改超管个人信息
	sysManageGroup, err := domain.GetSysManageGroup(orgId)
	if err != nil {
		log.Errorf("[UpdateOrgUser] GetSysManageGroup err:%v, orgId:%v", err, orgId)
		return errs.ManageGroupNotExist
	}
	sysUserIds := []int64{}
	if sysManageGroup.UserIds != nil {
		err = json.FromJson(*sysManageGroup.UserIds, &sysUserIds)
		if err != nil {
			log.Errorf("[UpdateOrgUser] json convert err:%v, orgId:%v", err, orgId)
			return errs.JSONConvertError
		}
	}
	if len(sysUserIds) > 0 {
		ok, _ := slice.Contain(sysUserIds, operator)
		ok2, _ := slice.Contain(sysUserIds, req.UserId)
		if !ok && ok2 {
			// 当前操作人不是超管，要修改的人是超管
			return errs.CannotEditSuperAdminInfo
		}
	}

	// 获取组织成员信息，包括非启用状态
	memberUserInfo, errSys := domain.GetOrgMemberInfoByUserId(orgId, req.UserId)
	//memberUserInfo, errSys := domain.GetOrgMemberBaseInfoListByUser(orgId, req.UserId)
	if errSys != nil {
		log.Errorf("[UpdateOrgUser] GetOrgMemberInfoByUserId err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
		return errSys
	}

	//globalUser, errSys := domain.GetGlobalUserByUserId(req.UserId)
	//if errSys != nil {
	//	log.Errorf("[UpdateOrgUser] GetGlobalUserByUserId err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
	//	return errSys
	//}

	// 获取管理组信息
	//managerGroup, errSys := domain.GetManageGroupListByUser(orgId, req.UserId)
	//if errSys != nil {
	//	log.Errorf("[UpdateOrgUser] GetManageGroupListByUser err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
	//	return errSys
	//}

	userUpd := mysql.Upd{}
	orgUserUpd := mysql.Upd{}

	// 修改姓名
	if req.Name != "" && req.Name != memberUserInfo.Name {
		userUpd[consts.TcName] = req.Name
	}

	// 判断状态是否能修改
	if req.Status != 0 && memberUserInfo.UserStatus != req.Status {
		// 自己不能改自己状态
		if req.UserId == operator {
			return errs.CannotChangeSelfStatus
		}
		orgUserUpd[consts.TcStatus] = req.Status
		orgUserUpd[consts.TcStatusChangerId] = operator
		orgUserUpd[consts.TcUpdator] = operator
		orgUserUpd[consts.TcStatusChangeTime] = time.Now()
		//errSys = domain.ModifyOrgMemberStatus(orgId, []int64{req.UserId}, req.Status, operator)
		//if errSys != nil {
		//	log.Errorf("[UpdateOrgUser] ModifyOrgMemberStatus err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
		//	return errSys
		//}
	}

	// 修改密码
	if req.Password != "" {
		errSys = RetrievePassword(orgvo.RetrievePasswordReqVo{
			OrgId: orgId,
			Input: vo.RetrievePasswordReq{
				Username:    memberUserInfo.LoginName,
				NewPassword: req.Password,
			}})
		if errSys != nil {
			log.Errorf("[UpdateOrgUser] RetrievePassword err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
			return errSys
		}
	}

	if len(userUpd) > 0 || len(orgUserUpd) > 0 {
		errSys = domain.UpdateOrgMemberInfo(orgId, req.UserId, userUpd, orgUserUpd)
		if errSys != nil {
			log.Errorf("[UpdateOrgUser] UpdateOrgMemberInfo err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
			return errSys
		}
	}

	// 更新用户的部门信息
	if req.DepartmentIds != nil {
		errSys = domain.ChangeUserDept(orgId, operator, req.UserId, req.DepartmentIds)
		if errSys != nil {
			log.Errorf("[UpdateOrgUser] ChangeUserDept err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
			return errSys
		}
	}

	// 更新用户的角色管理组信息
	if req.RoleGroupIds != nil {
		errSys = domain.ChangeUserManageGroup(orgId, operator, orgvo.ChangeUserManageGroupReq{
			UserId:         req.UserId,
			ManageGroupIds: req.RoleGroupIds,
		})
		if errSys != nil {
			log.Errorf("[UpdateOrgUser] GetManageGroupListByUser err:%v, orgId:%v, userId:%v", errSys, orgId, req.UserId)
			return errSys
		}
	}
	// 清除用户缓存信息
	cacheErr := domain.ClearBaseUserInfo(orgId, req.UserId)
	if cacheErr != nil {
		log.Errorf("[UpdateOrgUser] ClearBaseUserInfo err:%v, orgId:%v, userId:%v", cacheErr, orgId, req.UserId)
		return cacheErr
	}
	return nil
}

func GetOrgSuperAdminInfo(orgId int64) ([]*orgvo.GetOrgSuperAdminInfoData, errs.SystemErrorInfo) {
	return domain.GetOrgSuperAdminInfo(orgId)
}

func UpdateUserToSysManageGroup(orgId int64, userIds []int64, updateType int) errs.SystemErrorInfo {
	return domain.UpdateUserToSysManageGroup(orgId, userIds, updateType)
}
