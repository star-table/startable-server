package domain

import (
	"reflect"
	"strconv"
	"time"

	consts2 "gitea.bjx.cloud/allstar/platform-sdk/consts"

	"github.com/star-table/startable-server/app/facade/tablefacade"

	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	fsvo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	sconsts "github.com/star-table/startable-server/app/service/projectsvc/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/dao"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/id/snowflake"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// CreateMainChat 创建主群聊（创建项目时主动创建的群聊）的关联关系
func CreateMainChat(orgId, currentUserId int64, projectId int64, chatId string, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	var dbErr error
	//创建主群聊关联关系
	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	pos := po.PpmProProjectRelation{
		Id:           id,
		OrgId:        orgId,
		ProjectId:    projectId,
		RelationType: consts.IssueRelationTypeMainChat,
		RelationCode: chatId,
		Creator:      currentUserId,
		CreateTime:   time.Now(),
		IsDelete:     consts.AppIsNoDelete,
		Status:       consts.ProjectMemberEffective,
		Updator:      currentUserId,
		UpdateTime:   time.Now(),
		Version:      1,
	}
	if tx != nil {
		dbErr = mysql.TransInsert(tx, &pos)
	} else {
		dbErr = dao.InsertProjectRelation(pos)
	}
	if dbErr != nil {
		log.Errorf("[CreateMainChat] err: %v, project: %d", dbErr, projectId)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return id, nil
}

func CreateChat(orgId, currentUserId int64, projectId int64, chatIds []string) (int64, errs.SystemErrorInfo) {
	//创建关联关系
	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(chatIds))
	if err != nil {
		log.Error(err)
		return 0, err
	}
	pos := []po.PpmProProjectRelation{}
	for i, chatId := range chatIds {
		pos = append(pos, po.PpmProProjectRelation{
			Id:           ids.Ids[i].Id,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationType: consts.IssueRelationTypeChat,
			RelationCode: chatId,
			Creator:      currentUserId,
			CreateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
			Status:       consts.ProjectMemberEffective,
			Updator:      currentUserId,
			UpdateTime:   time.Now(),
			Version:      1,
		})
	}
	err1 := dao.InsertProjectRelationBatch(pos)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.MysqlOperateError
	}

	return int64(len(chatIds)), nil
}

func GetMainChatIdByProjectId(orgId, projectId int64) (string, errs.SystemErrorInfo) {
	po := po.PpmProProjectRelation{}
	err := mysql.SelectOneByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.IssueRelationTypeMainChat,
	}, &po)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return "", nil
		}
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return po.RelationCode, nil
}

// GetDefaultGroupChatSetting 群聊交互调整后，创建群聊，需要将群聊默认一个配置，存入配置表（ppm_pro_project_chat 表）中
func GetDefaultGroupChatSetting() *bo.GetFsTableSettingOfProjectChatBo {
	bo1 := &bo.GetFsTableSettingOfProjectChatBo{
		FsGroupChatSettingItems: bo.FsGroupChatSettingItems{
			CreateIssue:        1,
			CreateIssueComment: 1,
			UpdateIssueCase:    1,
			ModifyColumnsOfSend: []string{
				consts.GroupChatAllIssueColumnFlag,
			},
		},
	}

	return bo1
}

// CreateChatNew 新增群聊信息关联/新增群聊配置
// 关联后，会有一个默认配置，如：开启任务创建更新通知等
func CreateChatNew(orgId, currentUserId int64, projectId int64, chatIds []string, chatType int) (int64, errs.SystemErrorInfo) {
	// 查询项目下的表
	proAppIds, err := GetProjectAppIdsByProjectIds(orgId, []int64{projectId})
	if err != nil {
		log.Errorf("[CreateChatNew] projectId: %d, GetProjectAppIdsByProjectIds err: %v", projectId, err)
		return 0, nil
	}
	curProAppId := int64(0)
	if len(proAppIds) == 1 {
		curProAppId = proAppIds[0]
	}
	tablesMap, err := GetTableListMapByProAppIds(orgId, proAppIds)
	if err != nil {
		log.Errorf("[CreateChatNew] projectId: %d, GetTableListMapByProAppIds err: %v", projectId, err)
		return 0, nil
	}
	tables, ok := tablesMap[curProAppId]
	if !ok {
		err := errs.TableNotExist
		log.Errorf("[CreateChatNew] projectId: %d, err: %v", projectId, err)
		return 0, nil
	}
	//创建群聊配置
	chatSettingJson := json.ToJsonIgnoreError(GetDefaultGroupChatSetting())
	pos := make([]po.PpmProProjectChat, 0)
	for _, table := range tables {
		for _, chatId := range chatIds {
			pos = append(pos, po.PpmProProjectChat{
				Id:           snowflake.Id(),
				OrgId:        orgId,
				ProjectId:    projectId,
				TableId:      table.TableId,
				ChatId:       chatId,
				ChatType:     chatType,
				ChatSettings: chatSettingJson,
				IsEnable:     consts.AppIsEnable,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
				IsDelete:     consts.AppIsNoDelete,
			})

			// 如果是主动创建的群聊，则创建额外的关联关系
			if chatType == consts.ChatTypeMain {
				if _, err := CreateMainChat(orgId, currentUserId, projectId, chatId, nil); err != nil {
					log.Errorf("[CreateChatNew] CreateMainChat err: %v, projectId: %d", err, projectId)
					return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}
		}
	}
	err1 := dao.InsertChatSettingBatch(pos)
	if err1 != nil {
		log.Errorf("[CreateChatNew] InsertChatSettingBatch err: %v, projectId: %d", err1, projectId)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	return int64(len(chatIds)), nil
}

func ChatInfo(chatIds []string, projectId int64, orgId int64) (*[]po.PpmProProjectRelation, errs.SystemErrorInfo) {
	po := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationCode: db.In(chatIds),
		consts.TcRelationType: consts.IssueRelationTypeChat,
	}, po)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	return po, nil
}

func GetMainChatIdMapByProjectIds(orgId int64, projectIds []int64) (map[int64]string, errs.SystemErrorInfo) {
	resMap := make(map[int64]string, 0)
	poArr := make([]po.PpmProProjectRelation, 0, len(projectIds))
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    db.In(projectIds),
		consts.TcRelationType: consts.IssueRelationTypeMainChat,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, &poArr)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return resMap, nil
		}
		log.Errorf("[GetMainChatIdByProjectIds] err: %v, projectIds: %v", err, json.ToJsonIgnoreError(projectIds))
		return resMap, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, relation := range poArr {
		resMap[relation.ProjectId] = relation.RelationCode
	}

	return resMap, nil
}

func CheckChatIdIsMainChat(orgId int64, projectId int64, chatId string) (bool, errs.SystemErrorInfo) {
	chatIdMap, err := GetMainChatIdMapByProjectIds(orgId, []int64{projectId})
	if err != nil {
		log.Errorf("[CheckChatIdIsMainChat] err: %v, projectId: %d", err, projectId)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if tmpChatId, ok := chatIdMap[projectId]; ok {
		if chatId == tmpChatId {
			return true, nil
		}
	}

	return false, nil
}

// GetProjectIdsInChatIsMain 在群聊中，查询一批项目中，该群聊是否是主动创建的关系，返回对应的 projectId 集合
func GetProjectIdsInChatIsMain(orgId int64, projectIds []int64, chatId string) ([]int64, errs.SystemErrorInfo) {
	// 理论上，一个项目只会有一个主动创建的群聊
	resProIds := make([]int64, 0, 1)
	if len(projectIds) < 1 {
		return resProIds, nil
	}
	chatIdMap, err := GetMainChatIdMapByProjectIds(orgId, projectIds)
	if err != nil {
		log.Errorf("[GetProjectIdsInChatIsMain] err: %v, chatId: %v", err, chatId)
		return resProIds, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, proId := range projectIds {
		if tmpChatId, ok := chatIdMap[proId]; ok && tmpChatId == chatId {
			resProIds = append(resProIds, proId)
		}
	}

	return resProIds, nil
}

// ChatInfoNew 配置独立到表后，查询已存在的群聊信息
func ChatInfoNew(orgId, projectId int64, chatIds []string) ([]po.ProjectChatObj, errs.SystemErrorInfo) {
	if len(chatIds) < 1 {
		return nil, nil
	}
	list, err := GetProjectChatList(orgId, projectId, db.Cond{
		consts.TcChatId: db.In(chatIds),
	})
	if err != nil {
		log.Errorf("[ChatInfoNew] projectId：%v, err: %v", projectId, err)
		return list, err
	}

	return list, nil
}

func GetProjectMainChatId(orgId, projectId int64) (string, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectMainChatId, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return "", err5
	}
	data, err := cache.Get(key)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if data != "" {
		return data, nil
	} else {
		po := po.PpmProProjectRelation{}
		err := mysql.SelectOneByCond(consts.TableProjectRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcProjectId:    projectId,
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcRelationType: consts.IssueRelationTypeMainChat,
		}, &po)
		if err != nil {
			if err == db.ErrNoMoreRows {
				return "", nil
			}
			log.Error(err)
			return "", errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		err = cache.SetEx(key, po.RelationCode, consts.GetCacheBaseExpire())
		if err != nil {
			return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return po.RelationCode, nil
	}
}

// GetProjectMainChatIdInChatTable 群聊配置独立的表中后的查询方式
func GetProjectMainChatIdInChatTable(orgId, projectId int64) (string, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectMainChatId, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return "", err5
	}

	data, err := cache.Get(key)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if data != "" {
		return data, nil
	} else {
		po := &po.PpmProProjectChat{}
		err := mysql.SelectOneByCond(consts.TableProjectChat, db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcChatType:  consts.ChatTypeMain,
			consts.TcProjectId: projectId,
		}, po)
		if err != nil {
			if err == db.ErrNoMoreRows {
				return "", nil
			} else {
				log.Errorf("[GetProjectMainChatIdInChatTable] projectId: %d, err: %v", projectId, err)
				return "", errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
		err = cache.SetEx(key, po.ChatId, consts.GetCacheBaseExpire())
		if err != nil {
			log.Errorf("[GetProjectMainChatIdInChatTable] projectId: %d, err: %v", projectId, err)
			return "", errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return po.ChatId, nil
	}
}

func ClearProjectMainChatCache(orgId, projectId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectMainChatId, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

// ClearSomeCache 清除缓存
func ClearSomeCache(keyTpl string, param map[string]interface{}) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(keyTpl, param)
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

// 增加/移除飞书群聊成员
func AddFsChatMembers(orgId int64, projectId int64, addIds []int64, delIds []int64, addDeptIds []int64, delDeptIds []int64) errs.SystemErrorInfo {
	if len(addIds) == 0 && len(delIds) == 0 && len(addDeptIds) == 0 && len(delDeptIds) == 0 {
		return nil
	}

	// 获取组织信息
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	// 不是飞书直接返回
	if baseOrgInfo.SourceChannel != consts2.SourceChannelFeishu {
		return nil
	}
	if baseOrgInfo.OutOrgId == "" {
		return errs.CannotBindChat
	}

	projectInfo, err := GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	addIds = append(addIds, projectInfo.Owner)
	chatId, err := GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if chatId == "" {
		log.Info("当前项目未创建飞书群聊")
		return nil
	}
	if len(addDeptIds) > 0 {
		if ok, _ := slice.Contain(addDeptIds, int64(0)); ok {
			userIdsResp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{OrgId: orgId})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return userIdsResp.Error()
			}
			addIds = userIdsResp.Data
		} else {
			userIdsResp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: addDeptIds,
			})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return userIdsResp.Error()
			}
			addIds = append(addIds, userIdsResp.Data.UserIds...)
		}
	}
	if len(delDeptIds) > 0 {
		if ok, _ := slice.Contain(delDeptIds, int64(0)); ok {
			userIdsResp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{OrgId: orgId})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return userIdsResp.Error()
			}
			delIds = userIdsResp.Data
		} else {
			userIdsResp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: delDeptIds,
			})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return userIdsResp.Error()
			}
			delIds = append(delIds, userIdsResp.Data.UserIds...)
		}
	}
	delIds, addIds = util.GetDifMemberIds(delIds, addIds)

	allUserIds := append(addIds, delIds...)
	//获取人员openId
	userInfo := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		UserIds: allUserIds,
		OrgId:   orgId,
	})
	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return userInfo.Error()
	}
	if len(userInfo.BaseUserInfos) == 0 {
		return nil
	}

	addOpenIds := []string{}
	delOpenIds := []string{}
	for _, respVo := range userInfo.BaseUserInfos {
		if ok, _ := slice.Contain(addIds, respVo.UserId); ok {
			addOpenIds = append(addOpenIds, respVo.OutUserId)
		} else {
			delOpenIds = append(delOpenIds, respVo.OutUserId)
		}
	}

	tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	//每次只能传200个（增加群聊成员）
	count := len(addOpenIds)
	for i := 0; i < count; i += 200 {
		max := i + 200
		if max > count {
			max = count
		}
		resp, err2 := tenant.AddChatUser(fsvo.UpdateChatMemberReqVo{
			OpenIds: addOpenIds[i:max],
			ChatId:  chatId,
		})
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
		if resp.Code != 0 {
			log.Error("增加群聊成员失败:" + resp.Msg)
		}
	}

	//每次只能传200个（移除群聊成员）
	delCount := len(delOpenIds)
	for i := 0; i < delCount; i += 200 {
		max := i + 200
		if max > delCount {
			max = delCount
		}
		resp, err2 := tenant.RemoveChatUser(fsvo.UpdateChatMemberReqVo{
			OpenIds: delOpenIds[i:max],
			ChatId:  chatId,
		})
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
		if resp.Code != 0 {
			log.Error("移除群聊成员失败:" + resp.Msg)
		}
	}

	return nil
}

// UpdateChatTitle 更新群聊名称
func UpdateChatTitle(orgId int64, projectId int64, title string) errs.SystemErrorInfo {
	chatId, err := GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if chatId == "" {
		log.Infof("[UpdateChatTitle] 当前项目未创建飞书群聊 projectId: %d", projectId)
		return nil
	}

	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if baseOrgInfo.OutOrgId == "" {
		return errs.CannotBindChat
	}
	tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	resp, updateErr := tenant.UpdateChat(fsvo.UpdateChatReqVo{
		ChatId: chatId,
		Name:   title,
	})
	if updateErr != nil {
		log.Error(updateErr)
		return errs.FeiShuOpenApiCallError
	}

	if resp.Code != 0 {
		log.Errorf("[UpdateChatTitle] 更新飞书群聊失败: %s, projectId: %d", resp.Msg, projectId)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	return nil
}

// GetProjectIdByMainChatId 查找该群聊对应的主项目
func GetProjectIdByMainChatId(orgId int64, chatId string) (int64, errs.SystemErrorInfo) {
	info := &po.PpmProProjectChat{}
	err := mysql.SelectOneByCond(consts.TableProjectChat, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcChatId:   chatId,
		consts.TcOrgId:    orgId,
		consts.TcChatType: consts.ChatTypeMain,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return 0, nil
		} else {
			log.Errorf("[GetProjectIdByMainChatId] orgId: %d, err: %v", orgId, err)
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	return info.ProjectId, nil
}

// GetProjectIdByChatId 查找群聊对应的推送配置，包含主群聊、非主群聊
func GetProjectIdByChatId(orgId int64, chatId string) ([]po.PpmProProjectChat, errs.SystemErrorInfo) {
	poArr := make([]po.PpmProProjectChat, 0)
	err := mysql.SelectAllByCond(consts.TableProjectChat, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcChatId:   chatId,
		consts.TcOrgId:    orgId,
	}, &poArr)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return poArr, nil
		} else {
			log.Errorf("[GetProjectIdByChatId] orgId: %d, err: %v", orgId, err)
			return poArr, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	return poArr, nil
}

func GetFsChatPushSettings(orgId int64, chatId string,
) (*bo.GetFsProjectChatPushSettingsForChatBo, errs.SystemErrorInfo) {
	// key, err5 := util.ParseCacheKey(sconsts.CacheFsPushSettings, map[string]interface{}{
	// 	consts.CacheKeyOrgIdConstName:     orgId,
	// 	consts.CacheKeyProjectIdConstName: projectId,
	// })
	// if err5 != nil {
	// 	log.Error(err5)
	// 	return nil, err5
	// }

	res := bo.GetFsProjectChatPushSettingsForChatBo{
		Tables:              make([]bo.GetFsProjectChatPushSettingsForChatBoTable, 0),
		ModifyColumnsOfSend: make([]string, 0),
	}
	chatSettingPoArr, err := GetProjectIdByChatId(orgId, chatId)
	if err != nil {
		log.Errorf("[GetProjectFsPushSettings] err: %v", err)
		return nil, err
	}
	chatSettingArr, err := TransferProjectChatSettingToChat(chatSettingPoArr)
	if err != nil {
		log.Errorf("[GetProjectFsPushSettings] err: %v", err)
		return nil, err
	}
	if len(chatSettingArr) > 0 {
		res = chatSettingArr[0]
	}

	return &res, nil
}

// GetProjectFsPushSettingsOfTable 获取项目下，单个表的群聊（项目自带的群聊）配置信息
func GetProjectFsPushSettingsOfTable(orgId int64, chatId string, projectIds,
	tableIds []int64, tx sqlbuilder.Tx,
) (*bo.GetFsTableSettingOfProjectChatBo, errs.SystemErrorInfo) {
	var dbErr error
	settingBo := bo.GetFsTableSettingOfProjectChatBo{}
	settingPo := po.PpmProProjectChat{}
	if len(tableIds) == 0 {
		return &settingBo, nil
	}
	cond := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcTableId:   db.In(tableIds),
		consts.TcProjectId: db.In(projectIds),
		consts.TcChatId:    chatId,
	}
	if tx == nil {
		dbErr = mysql.SelectOneByCond(consts.TableProjectChat, cond, &settingPo)
	} else {
		dbErr = mysql.TransSelectOneByCond(tx, consts.TableProjectChat, cond, &settingPo)
	}
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return &settingBo, nil
		}
		log.Errorf("[GetProjectFsPushSettingsOfTable] orgId: %d, err: %v", orgId, dbErr)
		return &settingBo, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	json.FromJson(settingPo.ChatSettings, &settingBo)
	settingBo.SettingId = settingPo.Id
	settingBo.TableIdStr = strconv.FormatInt(settingPo.TableId, 10)

	return &settingBo, nil
}

// GetGroupChatDynSettingDefaultVal 一些设置值存储，包含群聊中“项目动态设置”的默认值
func GetGroupChatDynSettingDefaultVal() *bo.GetFsProjectChatPushSettingsBo {
	return &bo.GetFsProjectChatPushSettingsBo{
		ProjectId:           0,
		IsProjectChatNative: 0,
		Tables:              make([]bo.GetFsTableSettingOfProjectChatBo, 0),
	}
}

func GetGroupChatDynSettingDefaultValOld() *bo.GetFsProjectChatPushSettingsOldBo {
	return &bo.GetFsProjectChatPushSettingsOldBo{
		CreateIssue:                  1, //默认开启
		UpdateIssueOwner:             2,
		UpdateIssueStatus:            2,
		UpdateIssueProjectObjectType: 2,
		UpdateIssueTitle:             2,
		UpdateIssueTime:              2,
		CreateIssueComment:           2,
		UploadNewAttachment:          2,
		ProPrivacyStatus:             consts.ProSetPrivacyDisable, // 项目，默认不开启隐私模式
	}
}

func ClearProjectFsPushSettingsCache(orgId, projectId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheFsPushSettings, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

// GetProChatSettingsByProIdForOrigin 项目的自带群聊配置只有一个
func GetProChatSettingsByProIdForOrigin(orgId int64, projectId int64) (*po.PpmProProjectChat, errs.SystemErrorInfo) {
	chatSettingArr, err := GetProChatSettingsByProId(orgId, projectId, 1)
	if err != nil {
		log.Errorf("[GetProChatSettingsByProIdForOrigin] err: %v", err)
		return nil, err
	}
	if len(chatSettingArr) > 0 {
		return &chatSettingArr[0], nil
	}

	return nil, nil
}

// GetProChatSettingsByProId 通过 projectId 查询项目下所有的群聊配置
// isOrigin 是否只获取项目自带群聊 1只获取主群聊配置；2只获取非主群聊配置；-1都包含
func GetProChatSettingsByProId(orgId int64, projectId int64, isOrigin int) ([]po.PpmProProjectChat,
	errs.SystemErrorInfo) {
	chatSettingArr, err := GetProChatSettingsBatchByProIds(orgId, []int64{projectId}, isOrigin)
	if err != nil {
		log.Errorf("[GetProChatSettingsByProId] err: %v", err)
		return nil, err
	}

	return chatSettingArr, nil
}

// GetProChatSettingBoArrByTableIds 项目下，根据 tableId 查询群聊配置
func GetProChatSettingBoArrByTableIds(orgId int64, projectId int64, tableIds []int64) ([]po.PpmProProjectChat,
	errs.SystemErrorInfo) {
	poList := make([]po.PpmProProjectChat, 0, len(tableIds))
	if len(tableIds) < 1 {
		return poList, nil
	}
	poList, err := GetProChatSettingByCond(orgId, db.Cond{
		consts.TcProjectId: projectId,
		consts.TcTableId:   db.In(tableIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		log.Errorf("[GetProChatSettingBoArrByTableIds] err: %v", err)
		return nil, err
	}

	return poList, nil
}

// GetProChatSettingsBatchByProIds 获取一批项目绑定的群聊设置
// isOrigin 是否只获取项目自带群聊 1只获取主群聊配置；2只获取非主群聊配置；-1都包含
func GetProChatSettingsBatchByProIds(orgId int64, projectIds []int64, isOrigin int) ([]po.PpmProProjectChat,
	errs.SystemErrorInfo) {
	list := make([]po.PpmProProjectChat, 0)
	if len(projectIds) < 1 {
		return list, nil
	}
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: db.In(projectIds),
	}
	if isOrigin == 1 {
		cond1[consts.TcChatType] = consts.ChatTypeMain
	} else if isOrigin == 2 {
		cond1[consts.TcChatType] = consts.ChatTypeOut
	}
	err := mysql.SelectAllByCond(consts.TableProjectChat, cond1, &list)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return list, nil
		}
		log.Errorf("[GetProChatSettingsBatchByProIds] orgId: %d, err: %v", orgId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return list, nil
}

func GetProChatSettingByProIdsAndTableIdsAndChatId(orgId int64, projectIds, tableIds []int64,
	chatId string) ([]po.PpmProProjectChat, errs.SystemErrorInfo) {
	poList := make([]po.PpmProProjectChat, 0, len(tableIds))
	if len(tableIds) < 1 {
		return poList, nil
	}
	poList, err := GetProChatSettingByCond(orgId, db.Cond{
		consts.TcProjectId: db.In(projectIds),
		consts.TcTableId:   db.In(tableIds),
		consts.TcOrgId:     orgId,
		consts.TcChatId:    chatId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		log.Errorf("[GetProChatSettingByProIdsAndTableIdsAndChatId] err: %v", err)
		return nil, err
	}

	return poList, nil
}

func GetProChatSettingByProIdsAndChatId(orgId int64, projectIds []int64,
	chatId string) ([]po.PpmProProjectChat, errs.SystemErrorInfo) {
	poList := make([]po.PpmProProjectChat, 0)
	if len(projectIds) < 1 {
		return poList, nil
	}
	poList, err := GetProChatSettingByCond(orgId, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: db.In(projectIds),
		consts.TcChatId:    chatId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		log.Errorf("[GetProChatSettingByProIdsAndChatId] err: %v", err)
		return nil, err
	}

	return poList, nil
}

func GetProChatSettingByChatId(orgId int64, chatId string) ([]po.PpmProProjectChat, errs.SystemErrorInfo) {
	poList := make([]po.PpmProProjectChat, 0)
	if len(chatId) < 1 {
		return poList, nil
	}
	poList, err := GetProChatSettingByCond(orgId, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcChatId:   chatId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		log.Errorf("[GetProChatSettingByChatId] err: %v", err)
		return nil, err
	}

	return poList, nil
}

func GetProChatSettingByCond(orgId int64, condition db.Cond) ([]po.PpmProProjectChat, errs.SystemErrorInfo) {
	list := make([]po.PpmProProjectChat, 0)
	if len(condition) < 1 {
		return list, nil
	}
	condition[consts.TcOrgId] = orgId
	condition[consts.TcIsDelete] = consts.AppIsNoDelete
	err := mysql.SelectAllByCond(consts.TableProjectChat, condition, &list)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return list, nil
		}
		log.Errorf("[GetProChatSettingByCond] orgId: %d, err: %v", orgId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return list, nil
}

// GetProTableChatSettingOfOriginChat 获取项目自带的群的表的群聊配置
func GetProTableChatSettingOfOriginChat(orgId int64, projectId int64, tableId int64) (po.PpmProProjectChat, errs.SystemErrorInfo) {
	chatSetting := po.PpmProProjectChat{}
	if tableId == 0 {
		return chatSetting, nil
	}
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcTableId:   tableId,
		consts.TcChatType:  consts.ChatTypeMain,
	}
	err := mysql.SelectAllByCond(consts.TableProjectChat, cond1, &chatSetting)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return chatSetting, nil
		}
		log.Errorf("[GetProTableChatSettingOfOriginChat] tableId: %d, err: %v", tableId, err)
		return chatSetting, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return chatSetting, nil
}

func GetProTableChatSettingBoByChatId(orgId int64, projectId int64, tableId int64, chatId string) (bo.GetFsTableSettingOfProjectChatBo, errs.SystemErrorInfo) {
	bo := bo.GetFsTableSettingOfProjectChatBo{}
	chatSettingPo, err := GetProTableChatSettingPoByChatId(orgId, projectId, tableId, chatId)
	if err != nil {
		log.Errorf("[GetProTableChatSettingBoByChatId]  err: %v", err)
		return bo, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bo = ConvertOneTableChatSettingToBo(chatSettingPo)

	return bo, nil
}

func ConvertOneTableChatSettingToBo(settingPo po.PpmProProjectChat) bo.GetFsTableSettingOfProjectChatBo {
	boData := bo.GetFsTableSettingOfProjectChatBo{}
	json.FromJson(settingPo.ChatSettings, &boData)
	boData.SettingId = settingPo.Id

	return boData
}

func ConvertOneTableChatSettingToBoBatch(settingPoArr []po.PpmProProjectChat) []bo.GetFsTableSettingOfProjectChatBo {
	boList := make([]bo.GetFsTableSettingOfProjectChatBo, 0, len(settingPoArr))
	for _, settingPo := range settingPoArr {
		boData := bo.GetFsTableSettingOfProjectChatBo{}
		json.FromJson(settingPo.ChatSettings, &boData)
		boData.SettingId = settingPo.Id

		boList = append(boList, boData)
	}

	return boList
}

func GetProTableChatSettingPoByChatId(orgId int64, projectId int64, tableId int64, chatId string) (po.PpmProProjectChat, errs.SystemErrorInfo) {
	chatSetting := po.PpmProProjectChat{}
	if tableId == 0 {
		return chatSetting, nil
	}
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcTableId:   tableId,
		consts.TcChatId:    chatId,
	}
	err := mysql.SelectAllByCond(consts.TableProjectChat, cond1, &chatSetting)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return chatSetting, nil
		}
		log.Errorf("[GetProTableChatSettingByChatId] tableId: %d, err: %v", tableId, err)
		return chatSetting, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return chatSetting, nil
}

func GetProTableChatSettingBoArrByTableId(orgId int64, projectId int64, tableId int64) ([]bo.GetFsTableSettingOfProjectChatBo, errs.SystemErrorInfo) {
	result := make([]bo.GetFsTableSettingOfProjectChatBo, 0)
	if tableId == 0 {
		return result, nil
	}
	poList := make([]po.PpmProProjectChat, 0)
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcTableId:   tableId,
	}
	err := mysql.SelectAllByCond(consts.TableProjectChat, cond1, &poList)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		log.Errorf("[GetProTableChatSettingByChatId] tableId: %d, err: %v", tableId, err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	proSettingBo, err := TransferProjectChatSettingBatch(poList)
	if err != nil {
		log.Errorf("[GetProTableChatSettingBoArrByTableId] tableId: %d, err: %v", tableId, err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(proSettingBo) > 0 {
		result = proSettingBo[0].Tables
	}

	return result, nil
}

// GetProTableChatSettingBoArrByProId 通过项目获取群聊配置，每个群，取该项目中随机一个表关联的配置
func GetProTableChatSettingBoArrByProId(orgId int64, projectId int64) ([]bo.GetFsTableSettingOfProjectChatBo, errs.SystemErrorInfo) {
	result := make([]bo.GetFsTableSettingOfProjectChatBo, 0)
	if projectId == 0 {
		return result, nil
	}
	poList := make([]po.PpmProProjectChat, 0)
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
	}
	err := mysql.SelectAllByCond(consts.TableProjectChat, cond1, &poList)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		log.Errorf("[GetProTableChatSettingBoArrByProId] projectId: %d, err: %v", projectId, err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	proSettingBoArr, err := TransferProjectChatSettingBatch(poList)
	if err != nil {
		log.Errorf("[GetProTableChatSettingBoArrByProId] projectId: %d, err: %v", projectId, err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	// 按 chatId 分组
	group := map[string][]bo.GetFsTableSettingOfProjectChatBo{}
	for _, proSettingBo := range proSettingBoArr {
		for _, tableSetting := range proSettingBo.Tables {
			if _, ok := group[tableSetting.ChatId]; ok {
				group[tableSetting.ChatId] = []bo.GetFsTableSettingOfProjectChatBo{
					tableSetting,
				}
			} else {
				group[tableSetting.ChatId] = append(group[tableSetting.ChatId], tableSetting)
			}
		}
	}
	for _, settingBo := range group {
		// eg: 应用模板，应用后，快速关联另一个群，粒度只能到项目，无法精确到表，所以直接取第 0 个
		result = append(result, settingBo[0])
	}

	return result, nil
}

// GetProTableChatSetting 获取表的的群聊设置
func GetProTableChatSetting(orgId int64, projectId int64, tableId int64) ([]po.PpmProProjectChat, errs.SystemErrorInfo) {
	list := make([]po.PpmProProjectChat, 0)
	if tableId == 0 {
		return list, nil
	}
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcTableId:   tableId,
	}
	err := mysql.SelectAllByCond(consts.TableProjectChat, cond1, &list)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return list, nil
		}
		log.Errorf("[GetProTableChatSetting] orgId: %d, err: %v", orgId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return list, nil
}

// TransferChatSettingForOneProject 一个项目下的群聊配置转换为 bo 对象
func TransferChatSettingForOneProject(proChatSettingArr []po.PpmProProjectChat) (*bo.GetFsProjectChatPushSettingsBo, errs.SystemErrorInfo) {
	projectChatBoArr, err := TransferProjectChatSettingBatch(proChatSettingArr)
	if err != nil {
		log.Errorf("[TransferChatSettingForOneProject]  err: %v", err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(projectChatBoArr) > 0 {
		return &projectChatBoArr[0], nil
	}

	return nil, nil
}

// TransferProjectChatSettingBatch 查询到的群聊配置转为前端可识别的配置对象，针对项目的结构
func TransferProjectChatSettingBatch(proChatSettingArr []po.PpmProProjectChat) ([]bo.GetFsProjectChatPushSettingsBo, errs.SystemErrorInfo) {
	configBoArr := make([]bo.GetFsProjectChatPushSettingsBo, 0, len(proChatSettingArr))
	// 按项目 id 分组
	configBoGroupByProId := make(map[int64][]po.PpmProProjectChat, 0)
	for _, item := range proChatSettingArr {
		configBoGroupByProId[item.ProjectId] = append(configBoGroupByProId[item.ProjectId], item)
	}

	for proId, settingBoArr := range configBoGroupByProId {
		tmpProSettings := make([]bo.GetFsTableSettingOfProjectChatBo, 0, len(settingBoArr))
		for _, settingBo := range settingBoArr {
			tablesConfig := bo.GetFsTableSettingOfProjectChatBo{}
			settingObj := bo.FsGroupChatSettingItems{}
			if settingBo.ChatSettings != "" {
				json.FromJson(settingBo.ChatSettings, &settingObj)
			}
			tablesConfig.SettingId = settingBo.Id
			tablesConfig.ChatId = settingBo.ChatId

			tablesConfig.CreateIssue = settingObj.CreateIssue
			tablesConfig.UpdateIssueCase = settingObj.UpdateIssueCase
			tablesConfig.CreateIssueComment = settingObj.CreateIssueComment
			tablesConfig.ModifyColumnsOfSend = settingObj.ModifyColumnsOfSend

			tmpProSettings = append(tmpProSettings, tablesConfig)
		}
		configBoArr = append(configBoArr, bo.GetFsProjectChatPushSettingsBo{
			ProjectId:           proId,
			IsProjectChatNative: settingBoArr[0].ChatType,
			Tables:              tmpProSettings,
		})
	}

	return configBoArr, nil
}

// TransferProjectChatSettingToChat 将配置数据转换为针对群聊的推送配置数据结构
// proChatSettingArr 必须是同一个群聊下的所有配置
func TransferProjectChatSettingToChat(proChatSettingArr []po.PpmProProjectChat) ([]bo.
	GetFsProjectChatPushSettingsForChatBo, errs.SystemErrorInfo,
) {
	configBoArr := make([]bo.GetFsProjectChatPushSettingsForChatBo, 0)
	// 按 chat id 分组
	configBoGroupByChatId := make(map[string][]po.PpmProProjectChat, 0)
	for _, item := range proChatSettingArr {
		if _, ok := configBoGroupByChatId[item.ChatId]; !ok {
			configBoGroupByChatId[item.ChatId] = []po.PpmProProjectChat{
				item,
			}
		} else {
			configBoGroupByChatId[item.ChatId] = append(configBoGroupByChatId[item.ChatId], item)
		}
	}

	for _, settingBoArr := range configBoGroupByChatId {
		tmpTables := make([]bo.GetFsProjectChatPushSettingsForChatBoTable, 0, len(settingBoArr))
		tmpSettings := po.PpmProProjectChat{}
		tmpConfigBo := bo.GetFsProjectChatPushSettingsForChatBo{}
		for _, settingBo := range settingBoArr {
			tableInfo := bo.GetFsProjectChatPushSettingsForChatBoTable{}
			tableInfo.ProjectId = settingBo.ProjectId
			tableInfo.TableId = strconv.FormatInt(settingBo.TableId, 10)
			tmpTables = append(tmpTables, tableInfo)
			tmpSettings = settingBo
		}
		if tmpSettings.ChatSettings != "" {
			json.FromJson(tmpSettings.ChatSettings, &tmpConfigBo)
		}
		tmpConfigBo.Tables = tmpTables
		configBoArr = append(configBoArr, tmpConfigBo)
	}

	return configBoArr, nil
}

func GetProjectChatListOfNotMain(orgId, projectId int64, chatType int) ([]po.ProjectChatObj, errs.SystemErrorInfo) {
	return GetProjectChatList(orgId, projectId, db.Cond{
		consts.TcChatType: chatType,
	})
}

// GetProjectChatList 获取项目的群聊列表（非主群聊）
// chatType 1主群聊；2非主群聊
// projectId 为 -1 表示忽略这个查询条件
func GetProjectChatList(orgId, projectId int64, otherCond db.Cond) ([]po.ProjectChatObj, errs.SystemErrorInfo) {
	list := make([]po.ProjectChatObj, 0)
	if projectId == 0 {
		return list, nil
	}
	conn, oriErr := mysql.GetConnect()
	if oriErr != nil {
		log.Errorf("[GetProjectChatList] projectId: %d, err: %v", projectId, oriErr)
		return list, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	cond1 := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if projectId != -1 {
		cond1[consts.TcProjectId] = projectId
	}
	if len(otherCond) > 0 {
		for columnName, val := range otherCond {
			cond1[columnName] = val
		}
	}
	oriErr = conn.Select(db.Raw("id, org_id, project_id, chat_id, chat_type")).From(consts.TableProjectChat).
		Where(cond1).GroupBy("project_id").All(&list)
	if oriErr != nil {
		if oriErr == db.ErrNoMoreRows {
			return list, nil
		}
		log.Errorf("[GetProjectChatList] orgId: %d, err: %v", orgId, oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}

	return list, nil
}

// GetProjectIdsByChatId 根据群聊 id，查询与其关联的项目 id 列表
func GetProjectIdsByChatId(orgId int64, chatId string) ([]int64, errs.SystemErrorInfo) {
	poArr, err := GetProChatObjsGroupByProjectIdByChatId(orgId, chatId)
	if err != nil {
		log.Errorf("[GetProjectIdsByChatId] err: %v, orgId: %d, chatId: %s", err, orgId, chatId)
		return nil, err
	}
	projectIds := make([]int64, 0)
	for _, chatConfigObj := range poArr {
		projectIds = append(projectIds, chatConfigObj.ProjectId)
	}

	return projectIds, nil
}

// GetProChatObjsGroupByProjectIdByChatId 获取项目关联的群聊配置，按项目 id 分组
func GetProChatObjsGroupByProjectIdByChatId(orgId int64, chatId string) ([]po.ProjectChatObj, errs.SystemErrorInfo) {
	list := make([]po.ProjectChatObj, 0)
	if chatId == "" {
		return list, nil
	}
	list, err := GetProjectChatList(orgId, -1, db.Cond{
		consts.TcChatId: chatId,
	})
	if err != nil {
		log.Errorf("[GetProChatObjsGroupByProjectIdByChatId] err: %v, orgId: %d, chatId: %s", err, orgId, chatId)
		return nil, err
	}

	return list, nil
}

// CheckChatSettingIsForAllTable 检查用户的配置是否是针对项目下所有的表，用于新建表时检查
func CheckChatSettingIsForAllTable(orgId int64, projectId int64) (bool, errs.SystemErrorInfo) {
	// 查询项目下的表
	tablesMap, err := GetTableListMapByProAppIds(orgId, []int64{projectId})
	if err != nil {
		log.Errorf("[CheckChatSettingIsForAllTable] projectId: %d, GetTableListMapByProAppIds err: %v", projectId, err)
		return false, err
	}
	tables, ok := tablesMap[projectId]
	if !ok {
		err := errs.TableNotExist
		log.Errorf("[CheckChatSettingIsForAllTable] projectId: %d, err: %v", projectId, err)
		return false, err
	}
	tableIds := make([]int64, 0)
	for _, item := range tables {
		tableIds = append(tableIds, item.TableId)
	}
	// 查询表对应的群聊推送配置数据
	settingPoArr, err := GetProChatSettingBoArrByTableIds(orgId, projectId, tableIds)
	if err != nil {
		log.Errorf("[CheckChatSettingIsForAllTable] projectId: %d, err: %v", projectId, err)
		return false, err
	}
	// 对于新增的表，还没有配置，所以会少一条配置
	if len(settingPoArr) == (len(tableIds) - 1) {
		return true, nil
	} else {
		return false, nil
	}
}

// CreateGroupChatSettingForNewTable 为新 table 增加一个群聊配置（被哪些群聊绑定了，则对应群聊都需要加上该 table 的配置）
func CreateGroupChatSettingForNewTable(orgId, projectId, tableId int64, opUserId int64) errs.SystemErrorInfo {
	// 查询哪些群绑定了该项目（projectId）
	settingPoArr, err := GetProChatSettingsByProId(orgId, projectId, -1)
	if err != nil {
		log.Errorf("[CreateGroupChatSettingForOneTable] err: %v, projectId: %d", err, projectId)
		return err
	}
	chatIds := make([]string, 0, len(settingPoArr))
	// 按 chatId 分组
	settingMapByChatId := make(map[string][]po.PpmProProjectChat, len(settingPoArr))
	for _, setting := range settingPoArr {
		chatIds = append(chatIds, setting.ChatId)
		if _, ok := settingMapByChatId[setting.ChatId]; ok {
			settingMapByChatId[setting.ChatId] = append(settingMapByChatId[setting.ChatId], setting)
		} else {
			settingMapByChatId[setting.ChatId] = []po.PpmProProjectChat{
				setting,
			}
		}
	}
	chatIds = slice.SliceUniqueString(chatIds)
	// 查询项目下的表
	proAppIds, err := GetProjectAppIdsByProjectIds(orgId, []int64{projectId})
	if err != nil {
		log.Errorf("[CreateGroupChatSettingForOneTable] projectId: %d, GetProjectAppIdsByProjectIds err: %v", projectId, err)
		return err
	}
	curProAppId := int64(0)
	if len(proAppIds) > 0 {
		curProAppId = proAppIds[0]
	}
	tablesMap, err := GetTableListMapByProAppIds(orgId, []int64{curProAppId})
	if err != nil {
		log.Errorf("[CreateGroupChatSettingForOneTable] projectId: %d, GetTableListMapByProAppIds err: %v", projectId, err)
		return err
	}
	tables, ok := tablesMap[curProAppId]
	if !ok {
		err := errs.TableNotExist
		log.Errorf("[CreateGroupChatSettingForOneTable] projectId: %d, err: %v", projectId, err)
		return err
	}
	tableIds := make([]int64, 0)
	for _, item := range tables {
		tableIds = append(tableIds, item.TableId)
	}

	// 检查这些群中，该项目是否是选择了项目下所有的表
	needAddSettingChatIds := make([]string, 0)
	for chatId, settingArr := range settingMapByChatId {
		// 因为新增了表，所以已有的表配置一定比表数量少 1
		if len(settingArr) == len(tableIds)-1 {
			needAddSettingChatIds = append(needAddSettingChatIds, chatId)
		}
	}
	// 满足条件（选择了项目下所有的表），则创建一个新的配置记录
	pos := make([]po.PpmProProjectChat, 0, len(needAddSettingChatIds))
	for _, chatId := range needAddSettingChatIds {
		// 因为从前面遍历出来的，所以一定能匹配上
		tmpSettings, _ := settingMapByChatId[chatId]
		if len(tmpSettings) > 0 {
			pos = append(pos, po.PpmProProjectChat{
				Id:           snowflake.Id(),
				OrgId:        orgId,
				ProjectId:    projectId,
				TableId:      tableId,
				ChatId:       chatId,
				ChatType:     tmpSettings[0].ChatType,
				ChatSettings: tmpSettings[0].ChatSettings,
				IsEnable:     consts.AppIsEnable,
				Creator:      opUserId,
				CreateTime:   time.Now(),
				Updator:      opUserId,
				UpdateTime:   time.Now(),
				Version:      1,
				IsDelete:     consts.AppIsNoDelete,
			})
		}
	}
	if len(pos) > 0 {
		dbErr := dao.InsertChatSettingBatch(pos)
		if dbErr != nil {
			log.Errorf("[CreateGroupChatSettingForOneTable] err: %v, tableId: %d", dbErr, tableId)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
	}

	return nil
}

// DeleteGroupChatSettingForOneTable 删除一个表的群聊配置
func DeleteGroupChatSettingForOneTable(orgId, projectId, tableId int64, opUserId int64) errs.SystemErrorInfo {
	// 查询表在哪些群中
	cond := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcTableId:   tableId,
	}
	_, dbErr := mysql.UpdateSmartWithCond(consts.TableProjectChat, cond, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  opUserId,
	})
	if dbErr != nil {
		log.Errorf("[DeleteGroupChatSettingForOneTable] err: %v, tableId: %d", dbErr, tableId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

// DeleteGroupChatSettings 删除一个项目的某个群聊配置
func DeleteGroupChatSettings(orgId int64, projectId int64, chatId string) errs.SystemErrorInfo {
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		cond := db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcProjectId: projectId,
			consts.TcChatId:    chatId,
		}
		_, dbErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, cond, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if dbErr != nil {
			log.Error(dbErr)
			return dbErr
		}

		_, dbErr = mysql.TransUpdateSmartWithCond(tx, consts.TableProjectRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcProjectId:    projectId,
			consts.TcRelationCode: chatId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if dbErr != nil {
			log.Error(dbErr)
			return dbErr
		}
		return nil
	})

	if transErr != nil {
		log.Errorf("[DeleteGroupChatSettings] failed, orgId:%d, projectId:%d, err:%v", orgId, projectId, transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	return nil
}

// DeleteGroupChatSettingsByProjectId 删除一个项目的群聊配置
func DeleteGroupChatSettingsByProjectId(orgId int64, projectId int64) errs.SystemErrorInfo {
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		cond := db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcProjectId: projectId,
		}
		_, dbErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, cond, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if dbErr != nil {
			log.Errorf("[DeleteGroupChatSettingsByProjectId] err: %v, projectId: %d", dbErr, projectId)
			return dbErr
		}

		return nil
	})
	if transErr != nil {
		log.Errorf("[DeleteGroupChatSettingsByProjectId] failed, orgId:%d, projectId:%d, err:%v", orgId, projectId, transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	return nil
}

func DeleteFsProjectChatSettingByChatIdAndTableIds(tx sqlbuilder.Tx, orgId, userId int64, chatId string,
	tableIds []int64) errs.SystemErrorInfo {
	var dbErr error
	cond := db.Cond{
		consts.TcOrgId:   orgId,
		consts.TcChatId:  chatId,
		consts.TcTableId: db.In(tableIds),
	}
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	}
	if tx != nil {
		_, dbErr = mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, cond, upd)
	} else {
		_, dbErr = mysql.UpdateSmartWithCond(consts.TableProjectChat, cond, upd)
	}
	if dbErr != nil {
		log.Errorf("[DeleteFsProjectChatSettingByChatIdAndTableIds] err: %v, orgId: %d", dbErr, orgId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

func DeleteGroupSettingsByChatId(orgId int64, chatId string) errs.SystemErrorInfo {
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		_, dbErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, db.Cond{
			consts.TcOrgId:  orgId,
			consts.TcChatId: chatId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if dbErr != nil {
			log.Error(dbErr)
			return dbErr
		}

		_, dbErr = mysql.TransUpdateSmartWithCond(tx, consts.TableProjectRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcRelationCode: chatId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if dbErr != nil {
			log.Error(dbErr)
			return dbErr
		}
		return nil
	})

	if transErr != nil {
		log.Errorf("[DeleteGroupChatSettings] failed, orgId:%d, chatId:%s, err:%v", orgId, chatId, transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	return nil
}

// PushHelpCardAfterUpdateSetting 飞书群聊-保存了推送配置后，向群聊中推送一个帮助卡片 todo
// projectIds 群聊绑定的项目 id 列表
func PushHelpCardAfterUpdateSetting(orgId int64, chatId string, projectIds []int64) errs.SystemErrorInfo {
	orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if orgResp.Failure() {
		log.Errorf("[PushHelpCardAfterUpdateSetting] err: %v, orgId: %d", orgResp.Error(), orgId)
		return orgResp.Error()
	}
	// 查询是否绑定了多个项目
	bindProNum := len(projectIds)
	outOrgId := orgResp.BaseOrgInfo.OutOrgId
	cd := &pushV1.TemplateCard{}
	if bindProNum == 0 {
		cd = card.GetFsCardHelpOfBind0ProForGroupChat(orgId, chatId)
	} else if bindProNum == 1 {
		cd = card.GetFsCardHelpOfBind1ProForGroupChat()
	} else {
		cd = card.GetFsCardHelpOfBind2ProForGroupChat(bindProNum)
	}
	//if err := SendCardForGroupChat(outOrgId, chatId, cd); err != nil {
	//	log.Errorf("[PushHelpCardAfterUpdateSetting] err: %v, orgId: %d", err, orgId)
	//	return orgResp.Error()
	//}
	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      outOrgId,
		SourceChannel: orgResp.BaseOrgInfo.SourceChannel,
		ChatIds:       []string{chatId},
		CardMsg:       cd,
	})
	if errSys != nil {
		log.Errorf("[PushHelpCardAfterUpdateSetting] pushCard err:%v", errSys)
		return errSys
	}

	return nil
}

func DealTransForUpdateGroupChatSettings(orgId, userId int64, params *vo.UpdateFsProjectChatPushSettingsReq, tableIdsForUpdate, tableIdsForAdd, tableIdsForDel []int64,
	tableIdMapKeyByTableId map[int64]int64, proIdsForMain []int64,
) errs.SystemErrorInfo {
	newSetting := bo.FsGroupChatSettingItems{
		CreateIssue:         params.CreateIssue,
		CreateIssueComment:  params.CreateIssueComment,
		UpdateIssueCase:     params.UpdateIssueCase,
		ModifyColumnsOfSend: make([]string, 0),
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 处理更新的部分
		for _, tableId := range tableIdsForUpdate {
			tmpProjectId, ok := tableIdMapKeyByTableId[tableId]
			if !ok {
				continue
			}
			if err := UpdateFsProjectChatSettingForOneTable(orgId, userId, params.ChatID, tmpProjectId,
				strconv.FormatInt(tableId, 10), *params, tx); err != nil {
				log.Errorf("[UpdateFsProjectChatPushSettings] err: %v, orgId: %d, tableIdStr: %d", err, orgId,
					tableId)
				return err
			}
		}
		// 处理新增的部分
		insertData := make([]interface{}, 0, len(tableIdsForAdd))
		for _, tableId := range tableIdsForAdd {
			tmpProjectId, ok := tableIdMapKeyByTableId[tableId]
			if !ok {
				continue
			}
			chatType := consts.ChatTypeOut
			if exist, _ := slice.Contain(proIdsForMain, tmpProjectId); exist {
				chatType = consts.ChatTypeMain
			}
			insertData = append(insertData, po.PpmProProjectChat{
				Id:           snowflake.Id(),
				OrgId:        orgId,
				ProjectId:    tmpProjectId,
				TableId:      tableId,
				ChatId:       params.ChatID,
				ChatType:     chatType,
				ChatSettings: json.ToJsonIgnoreError(newSetting),
				IsEnable:     consts.AppIsEnable,
				Creator:      userId,
				CreateTime:   time.Now(),
				Updator:      userId,
				UpdateTime:   time.Now(),
				Version:      1,
				IsDelete:     consts.AppIsNoDelete,
			})
		}
		if len(insertData) > 0 {
			dbErr := mysql.TransBatchInsert(tx, &po.PpmProProjectChat{}, insertData)
			if dbErr != nil {
				log.Errorf("[DealTransForUpdateGroupChatSettings] TransBatchInsert err: %v, chatId: %s", dbErr,
					params.ChatID)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
			}
		}
		// 处理删除的部分
		if len(tableIdsForDel) > 0 {
			if err := DeleteFsProjectChatSettingByChatIdAndTableIds(tx, orgId, userId, params.ChatID,
				tableIdsForDel); err != nil {
				log.Errorf("[DealTransForUpdateGroupChatSettings] DeleteFsProjectChatSettingByChatIdAndTableIds err: %v, "+
					"orgId: %d, chatId: %s", err, orgId,
					params.ChatID)
				return err
			}
		}
		return nil
	})
	if transErr != nil {
		log.Errorf("[UpdateFsProjectChatPushSettings] trans err: %v, orgId: %d", transErr, orgId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	return nil
}

func UpdateFsProjectChatSettingForOneTable(orgId, userId int64, chatId string, projectId int64, tableIdStr string,
	params vo.UpdateFsProjectChatPushSettingsReq, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	var err errs.SystemErrorInfo
	// 初始设置
	tableId, oriErr := strconv.ParseInt(tableIdStr, 10, 64)
	if oriErr != nil {
		log.Errorf("[UpdateFsProjectChatSettingForOneTable] tableIdStr: %s, parse tableId err: %v", tableIdStr, oriErr)
		return err
	}
	tableSetting, err := GetProjectFsPushSettingsOfTable(orgId, chatId, []int64{projectId}, []int64{tableId}, tx)
	if err != nil {
		log.Errorf("[UpdateFsProjectChatSettingForOneTable] err: %v", err)
		return err
	}
	settingItemsBo := bo.FsGroupChatSettingItems{}
	copyer.Copy(tableSetting, &settingItemsBo)

	// 校验变更
	isUpdated := false
	newValue := settingItemsBo
	if ok, _ := slice.Contain([]int{1, 2}, params.CreateIssue); ok {
		if params.CreateIssue != tableSetting.CreateIssue {
			newValue.CreateIssue = params.CreateIssue
			isUpdated = true
		}
	}
	if ok, _ := slice.Contain([]int{1, 2}, params.CreateIssueComment); ok {
		if params.CreateIssueComment != tableSetting.CreateIssueComment {
			newValue.CreateIssueComment = params.CreateIssueComment
			isUpdated = true
		}
	}
	if ok, _ := slice.Contain([]int{1, 2}, params.UpdateIssueCase); ok {
		if params.UpdateIssueCase != tableSetting.UpdateIssueCase {
			newValue.UpdateIssueCase = params.UpdateIssueCase
			isUpdated = true
		}
	}
	if len(params.ModifyColumnsOfSend) < 1 {
		params.ModifyColumnsOfSend = []*string{
			str.ToPtr(consts.GroupChatAllIssueColumnFlag),
		}
	}
	if !reflect.DeepEqual(params.ModifyColumnsOfSend, tableSetting.ModifyColumnsOfSend) {
		copyer.Copy(params.ModifyColumnsOfSend, &newValue.ModifyColumnsOfSend)
		isUpdated = true
	}
	if !isUpdated {
		return nil
	}

	var dbErr error
	if tx != nil {
		_, dbErr = mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, db.Cond{
			consts.TcId: tableSetting.SettingId,
		}, mysql.Upd{
			consts.TcChatSetting: json.ToJsonIgnoreError(newValue),
			consts.TcUpdator:     userId,
		})
	} else {
		_, dbErr = mysql.UpdateSmartWithCond(consts.TableProjectChat, db.Cond{
			consts.TcId: tableSetting.SettingId,
		}, mysql.Upd{
			consts.TcChatSetting: json.ToJsonIgnoreError(newValue),
			consts.TcUpdator:     userId,
		})
	}
	if dbErr != nil {
		log.Errorf("[UpdateFsProjectChatSettingForOneTable] orgId: %d, chatSettingId: %d, err: %v",
			orgId, tableSetting.SettingId, dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

func ValidateAndFixInputForUpdateChatSetting(orgId int64, input *vo.UpdateFsProjectChatPushSettingsReq) errs.SystemErrorInfo {
	if input.ChatID == "" {
		if input.ProjectID > 0 {
			// 查询 projectId 对应的主群聊
			chatId, err := GetProjectMainChatId(orgId, input.ProjectID)
			if err != nil {
				log.Errorf("[ValidateAndFixInputForUpdateChatSetting] err: %v, chatId: %s, proId: %d", err, input.ChatID, input.ProjectID)
				return err
			}
			input.ChatID = chatId
		} else {
			return errs.ParamError
		}
	}
	if input.ChatID == "" {
		return errs.ParamError
	}

	return nil
}

// CheckIsAppCollaborator 检查当前用户是当前应用的协作人
func CheckIsAppCollaborator(orgId, appId int64, curUserId int64) bool {
	resp := tablefacade.CheckIsAppCollaborator(orgId, curUserId, appId, curUserId)
	if resp.Failure() {
		log.Errorf("[CheckIsAppCollaborator] err: %v, orgId: %d, appId: %d, userId: %d", resp.Failure(), orgId, appId, curUserId)
		return false
	}
	return resp.Data.Result
}

// CheckUserIsInFsChat 检查用户是否在群聊中
func CheckUserIsInFsChat(orgInfo *bo.BaseOrgInfoBo, userInfo *bo.BaseUserInfoBo, chatId string) (bool, errs.SystemErrorInfo) {
	outOrgId := orgInfo.OutOrgId
	if chatId == "" || outOrgId == "" {
		log.Infof("[CheckUserIsInFsChat] 参数不正确。outOrgId: %s, chatId: %s", outOrgId, chatId)
		return false, errs.ParamError
	}
	tenant, err := feishu.GetTenant(outOrgId)
	if err != nil {
		log.Errorf("[CheckUserIsInFsChat] GetTenant err: %v", err)
		return false, err
	}
	accessResp := orgfacade.GetFsAccessToken(orgvo.GetFsAccessTokenReqVo{
		OrgId:  orgInfo.OrgId,
		UserId: userInfo.UserId,
	})
	if accessResp.Failure() {
		log.Errorf("[CheckUserIsInFsChat] GetFsAccessToken err: %v, userId: %d", err, userInfo.UserId)
		return false, err
	}
	userAccessToken := accessResp.AccessToken

	resp, oriErr := tenant.CheckMemberIsInChat(userAccessToken, chatId)
	if oriErr != nil {
		log.Errorf("[CheckUserIsInFsChat] CheckMemberIsInChat err: %v", oriErr)
		return false, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, oriErr)
	}

	return resp.Data.IsInChat, nil
}

// AddUserToChat 将用户加入项目群聊中
func AddUserToChat(outOrgId, addOpenId, chatId string) errs.SystemErrorInfo {
	if chatId == "" || outOrgId == "" {
		log.Infof("[AddUserToChat] 参数不争正确。outOrgId: %s, chatId: %s", outOrgId, chatId)
		return nil
	}
	tenant, err := feishu.GetTenant(outOrgId)
	if err != nil {
		log.Errorf("[AddUserToProChat] GetTenant err: %v", err)
		return err
	}
	resp, oriErr := tenant.AddChatUser(fsvo.UpdateChatMemberReqVo{
		OpenIds: []string{addOpenId},
		ChatId:  chatId,
	})
	if oriErr != nil {
		log.Errorf("[AddUserToProChat] AddChatUser err: %v, chatId:%s", err, chatId)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
	}
	if resp.Code != 0 {
		log.Errorf("[AddUserToProChat] AddChatUser code err resp: %s, chatId: %s", json.ToJsonIgnoreError(resp), chatId)
		return errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, resp.Msg)
	}

	return nil
}
