package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/star-table/startable-server/common/model/vo/commonvo"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang/request"
	fsvo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	int642 "github.com/star-table/startable-server/common/core/util/slice/int64"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// 不同平台群聊
func AddChatFromPlatform(orgId, currentUserId int64, memberIds []int64, projectId int64, name string, remark *string,
	ownerIds []int64, departmentIds []int64, sourceChannel string) (string, errs.SystemErrorInfo) {
	switch sourceChannel {
	case sdk_const.SourceChannelFeishu:
		// 飞书创建群聊
		return AddFsChat(orgId, currentUserId, memberIds, projectId, name, remark, ownerIds, departmentIds)
	case sdk_const.SourceChannelDingTalk:
		// 钉钉创建群聊
		return AddDingChat(orgId, currentUserId, memberIds, projectId, name, ownerIds, departmentIds)
	default:
		return "", errs.SourceNotExist
	}
}

// AddFsChat 增加飞书群聊配置，后端创建的群聊，会给项目下所有的表创建配置
func AddFsChat(orgId, currentUserId int64, memberIds []int64, projectId int64, name string, remark *string, ownerIds []int64, departmentIds []int64) (string, errs.SystemErrorInfo) {
	//获取人员openId
	memberIds = append(memberIds, ownerIds...)
	//部门成员
	if len(departmentIds) > 0 {
		if ok, _ := slice.Contain(departmentIds, int64(0)); ok {
			userIdsResp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{OrgId: orgId})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return "", userIdsResp.Error()
			}
			memberIds = userIdsResp.Data
		} else {
			userIdsResp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: departmentIds,
			})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return "", userIdsResp.Error()
			}
			memberIds = append(memberIds, userIdsResp.Data.UserIds...)
		}
	}
	memberIds = slice.SliceUniqueInt64(memberIds)

	userInfo := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		UserIds: memberIds,
		OrgId:   orgId,
	})
	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return "", userInfo.Error()
	}
	openIds := []string{}
	//ownerOpenId := ""
	for _, respVo := range userInfo.BaseUserInfos {
		openIds = append(openIds, respVo.OutUserId)
		//if respVo.UserId == owner {
		//	ownerOpenId = respVo.OutUserId
		//}
	}
	//创建飞书群聊
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if baseOrgInfo.OutOrgId == "" {
		return "", errs.CannotBindChat
	}
	tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err1 != nil {
		log.Error(err1)
		return "", err1
	}
	trulyRemark := ""
	if remark != nil {
		trulyRemark = *remark
	}
	//查看有没有重复的名称（已删除的，因为创建的时候飞书的判断逻辑是所有信息和某个群一致，那么就返回旧有的群id，这里就判断名称相同，则在描述里面加点标识）
	isExist, isExistErr := mysql.IsExistByCond(consts.TableProject, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcName:  name,
		consts.TcId:    db.NotEq(projectId),
	})
	if isExistErr != nil {
		log.Error(isExistErr)
		return "", errs.MysqlOperateError
	}
	if isExist {
		trulyRemark += fmt.Sprintf("（项目ID:%d）", projectId)
	}

	batchNum := 200
	firstOpenIds := []string{}
	secondOpenIds := []string{}
	if len(openIds) <= batchNum {
		firstOpenIds = openIds
	} else {
		firstOpenIds = openIds[:batchNum]
		secondOpenIds = openIds[batchNum:]
	}
	resp, err2 := tenant.CreateChat(fsvo.CreateChatReqVo{
		Name:        name,
		Description: trulyRemark,
		OpenIds:     firstOpenIds,
	})
	if err2 != nil {
		log.Error(err2)
		return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err2)
	}
	if resp.Code != 0 {
		log.Errorf("[AddFsChat] 创建飞书群聊失败 resp: %s", json.ToJsonIgnoreError(resp))
		return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	//每次只能传200个（增加群聊成员）
	count := len(secondOpenIds)
	for i := 0; i < count; i += batchNum {
		max := i + batchNum
		if max > count {
			max = count
		}
		resp, err2 := tenant.AddChatUser(fsvo.UpdateChatMemberReqVo{
			OpenIds: secondOpenIds[i:max],
			ChatId:  resp.Data.ChatId,
		})
		if err2 != nil {
			log.Error(err2)
			return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
		if resp.Code != 0 {
			log.Error("增加群聊成员失败:" + resp.Msg)
			return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
	}

	////把群主转让给项目负责人
	//transferResp, err3 := tenant.UpdateChat(fsvo.UpdateChatReqVo{
	//	ChatId:      resp.NewData.ChatId,
	//	OwnerOpenId: ownerOpenId,
	//})
	//if err3 != nil {
	//	log.Error(err3)
	//	return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	//}
	//if transferResp.Code != 0 {
	//	log.Error("转让飞书群主失败:" + transferResp.Msg)
	//	return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	//}

	//创建主群聊关系
	_, insertErr := domain.CreateChatNew(orgId, currentUserId, projectId, []string{resp.Data.ChatId}, consts.ChatTypeMain)
	if insertErr != nil {
		log.Error(insertErr)
		return "", insertErr
	}
	// 不需要再加一条配置记录吧？
	//插入表
	//_, insertErr1 := domain.CreateChatNew(orgId, currentUserId, projectId, []string{resp.NewData.ChatId}, consts.ChatTypeOut)
	//if insertErr1 != nil {
	//	log.Error(insertErr1)
	//	return "", insertErr
	//}
	// 向群里发送帮助信息 start
	msgCard := card.GetFsCardGroupChatHelp(orgId, projectId, resp.Data.ChatId)

	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      baseOrgInfo.OutOrgId,
		SourceChannel: baseOrgInfo.SourceChannel,
		ChatIds:       []string{resp.Data.ChatId},
		CardMsg:       msgCard,
	})
	if errSys != nil {
		log.Errorf("[AddFsChat] pushCard err:%v", errSys)
		return "", errSys
	}
	//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
	//	ChatId:  resp.NewData.ChatId,
	//	MsgType: "interactive",
	//	Card:    msgCard,
	//})
	//if oriErr != nil {
	//	log.Error(oriErr)
	//	return "", err
	//}
	//if fsResp.Code != 0 {
	//	log.Error("AddFsChat 发送消息异常")
	//}
	// 向群里发送帮助信息 end

	return resp.Data.ChatId, nil
}

func AllChatList(userId int64, orgId int64, name *string) ([]fsvo.GroupData, errs.SystemErrorInfo) {
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if baseOrgInfo.OutOrgId == "" {
		return nil, errs.CannotBindChat
	}
	tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	result := []fsvo.GroupData{}
	pageToken := ""
	for {
		accessTokenResp := orgfacade.GetFsAccessToken(orgvo.GetFsAccessTokenReqVo{UserId: userId, OrgId: orgId})
		if accessTokenResp.Failure() {
			log.Error(accessTokenResp.Error())
			return nil, accessTokenResp.Error()
		}
		var resp *fsvo.GroupListRespVo
		var fsErr error
		if name == nil || *name == "" {
			resp, fsErr = tenant.GroupList(accessTokenResp.AccessToken, 0, pageToken)
		} else {
			resp, fsErr = tenant.ChatSearch(accessTokenResp.AccessToken, *name, 0, pageToken)
		}
		if fsErr != nil {
			return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
		if resp.Code != 0 {
			log.Error(resp.Msg)
			return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}

		if resp.Data != nil && resp.Data.Groups != nil && len(resp.Data.Groups) > 0 {
			result = append(result, resp.Data.Groups...)
		}

		if resp.Data.HasMore == false {
			break
		} else {
			pageToken = resp.Data.PageToken
		}
	}

	return result, nil
}

// ChatList 分页查询已关联项目的群聊
func ChatList(pageSize int64, currentUserId int64, projectId int64, orgId int64, lastRelationId *int64, name *string) (*vo.ChatListResp, errs.SystemErrorInfo) {
	result := vo.ChatListResp{
		Total: 0,
		List:  []*vo.ChatData{},
	}
	//查询已关联项目的群聊
	cond := db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.IssueRelationTypeChat,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}
	//if lastRelationId != nil && *lastRelationId != 0 {
	//	cond[consts.TcId] = db.Gt(*lastRelationId)
	//}
	data, err := domain.GetProjectRelationByCond(cond)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if data == nil {
		return &result, nil
	}

	relationChatIds := []string{}
	chatToRelationMap := map[string]int64{}
	needRelationChatIdsPage := []string{}
	for _, bo := range *data {
		relationChatIds = append(relationChatIds, bo.RelationCode)
		chatToRelationMap[bo.RelationCode] = bo.Id
		if lastRelationId != nil && *lastRelationId != 0 {
			if bo.Id > *lastRelationId {
				needRelationChatIdsPage = append(needRelationChatIdsPage, bo.RelationCode)
			}
		} else {
			needRelationChatIdsPage = append(needRelationChatIdsPage, bo.RelationCode)
		}
	}

	//获取项目主群聊
	mainChatId, err := domain.GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//取所有与自己相关的群聊
	allMyChatList, err := AllChatList(currentUserId, orgId, name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, groupData := range allMyChatList {
		if ok, _ := slice.Contain(relationChatIds, groupData.ChatId); ok {
			result.Total++
		}
	}
	//分页取数据
	var count int64 = 0
	for _, groupData := range allMyChatList {
		if count >= pageSize {
			break
		}
		if ok, _ := slice.Contain(needRelationChatIdsPage, groupData.ChatId); ok {
			chatData := vo.ChatData{
				Name:        groupData.Name,
				Description: &groupData.Description,
				OutChatID:   groupData.ChatId,
				Avatar:      groupData.Avatar,
				IsMain:      false,
			}
			if groupData.ChatId == mainChatId {
				chatData.IsMain = true
			}
			if _, ok := chatToRelationMap[groupData.ChatId]; ok {
				relationId := chatToRelationMap[groupData.ChatId]
				chatData.RelationID = &relationId
			}
			result.List = append(result.List, &chatData)
			count++
		}
	}

	return &result, nil
}

func reverse(arr *[]fsvo.GroupData) {
	var temp fsvo.GroupData
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

// UnrelatedChatList 分页查询未关联项目的群聊（非主群聊）
func UnrelatedChatList(pageSize int64, currentUserId int64, projectId int64, orgId int64, lastChatId, name *string) (*vo.ChatListResp, errs.SystemErrorInfo) {
	//查看项目信息
	projectInfo, projectErr := domain.GetProjectSimple(orgId, projectId)
	if projectErr != nil {
		log.Error(projectErr)
		return nil, errs.ProjectNotExist
	}
	result := vo.ChatListResp{
		Total: 0,
		List:  []*vo.ChatData{},
	}
	// 查询已关联项目的群聊
	data, err := domain.GetProjectChatListOfNotMain(orgId, projectId, consts.ChatTypeOut)
	if err != nil {
		log.Errorf("[UnrelatedChatList] projectId: %d, err: %v", projectId, err)
		return nil, err
	}
	if data == nil {
		return &result, nil
	}

	relationChatIds := []string{}
	for _, bo := range data {
		relationChatIds = append(relationChatIds, bo.ChatId)
	}
	//取所有与自己相关的群聊
	allMyChatList, err := AllChatList(currentUserId, orgId, name)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//倒序展示
	reverse(&allMyChatList)
	for _, groupData := range allMyChatList {
		if ok, _ := slice.Contain(relationChatIds, groupData.ChatId); !ok {
			result.Total++
		}
	}

	//取未展示的数据
	begin := 0
	for i, groupData := range allMyChatList {
		if lastChatId != nil && *lastChatId != "" && groupData.ChatId == *lastChatId {
			begin = i
			break
		}
	}
	if len(allMyChatList) <= begin+1 {
		return &result, nil
	}

	allMyChatList = allMyChatList[begin+1 : len(allMyChatList)]

	//获取项目主群聊
	mainChatId, err := domain.GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//分页取数据
	var count int64 = 0
	for _, groupData := range allMyChatList {
		if count >= pageSize {
			break
		}
		//未关联的群聊
		if ok, _ := slice.Contain(relationChatIds, groupData.ChatId); !ok {
			if groupData.ChatId == mainChatId && projectInfo.Owner != currentUserId {
				continue
			}
			chatData := vo.ChatData{
				Name:        groupData.Name,
				Description: &groupData.Description,
				OutChatID:   groupData.ChatId,
				Avatar:      groupData.Avatar,
				IsMain:      false,
			}
			if groupData.ChatId == mainChatId {
				chatData.IsMain = true
			}
			result.List = append(result.List, &chatData)
			count++
		}
	}

	return &result, nil
}

func AddChat(chatIds []string, orgId int64, currentUserId int64, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	//取所有与自己相关的群聊
	allMyChatList, err := AllChatList(currentUserId, orgId, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	allChatIds := []string{}
	for _, data := range allMyChatList {
		allChatIds = append(allChatIds, data.ChatId)
	}

	//防止重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	lockKey := consts.AddProjectChatLock + projectIdStr
	suc, cacheErr := cache.TryGetDistributedLock(lockKey, uid)
	if cacheErr != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, cacheErr)
		return nil, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	}

	//查看项目信息
	_, projectErr := domain.GetProjectSimple(orgId, projectId)
	if projectErr != nil {
		log.Error(projectErr)
		return nil, errs.ProjectNotExist
	}
	//查看群聊是否已关联
	relationData, infoErr := domain.ChatInfoNew(orgId, projectId, chatIds)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, infoErr
	}
	relatedChatIds := []string{}
	if relationData != nil {
		for _, relation := range relationData {
			relatedChatIds = append(relatedChatIds, relation.ChatId)
		}
	}

	//排除已关联的群聊
	needAddChatIds := []string{}
	if len(relatedChatIds) != 0 {
		for _, id := range chatIds {
			if ok, _ := slice.Contain(relatedChatIds, id); !ok {
				needAddChatIds = append(needAddChatIds, id)
			}
		}
	} else {
		needAddChatIds = chatIds
	}

	finalAddChatIds := []string{}
	for _, id := range needAddChatIds {
		if ok, _ := slice.Contain(allChatIds, id); ok {
			finalAddChatIds = append(finalAddChatIds, id)
		}
	}

	//关联群聊
	if len(finalAddChatIds) == 0 {
		return &vo.Void{ID: 0}, nil
	}
	id, insertErr := domain.CreateChatNew(orgId, currentUserId, projectId, finalAddChatIds, consts.ChatTypeOut)
	if insertErr != nil {
		log.Error(insertErr)
		return nil, insertErr
	}

	return &vo.Void{ID: id}, nil
}

// DisbandChat 解除群聊关联
func DisbandChat(chatIds []string, orgId int64, currentUserId int64, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	//查看群聊是否已关联
	info, infoErr := domain.ChatInfoNew(orgId, projectId, chatIds)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, infoErr
	}
	//获取项目主群聊
	mainChatId, err := domain.GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//项目主群聊只有项目负责人有权删除
	if mainChatId != "" {
		if ok, _ := slice.Contain(chatIds, mainChatId); ok {
			//查看项目信息
			projectInfo, projectErr := domain.GetProjectSimple(orgId, projectId)
			if projectErr != nil {
				log.Error(projectErr)
				return nil, errs.ProjectNotExist
			}
			if projectInfo.Owner != currentUserId {
				return nil, errs.CannotDisbandMainChat
			}
		}
	}
	needDisBandId := []int64{}
	if info != nil {
		for _, relation := range info {
			needDisBandId = append(needDisBandId, relation.Id)
		}
	}
	if len(needDisBandId) == 0 {
		return &vo.Void{ID: 0}, nil
	}

	// 解除群聊绑定
	count, deleteErr := mysql.UpdateSmartWithCond(consts.TableProjectChat, db.Cond{
		consts.TcId: db.In(needDisBandId),
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  currentUserId,
	})
	if deleteErr != nil {
		log.Errorf("[DisbandChat] err: %v", deleteErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, deleteErr)
	}

	return &vo.Void{ID: count}, nil
}

// FsChatDisbandCallback 删除群聊关联
func FsChatDisbandCallback(orgId int64, chatId string) errs.SystemErrorInfo {
	//查找该群聊对应的主群聊
	projectId, err := domain.GetProjectIdByMainChatId(orgId, chatId)
	if err != nil {
		log.Error(err)
		return err
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除所有与该群相关的项目关联
		_, deleteErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcChatId:   chatId,
			consts.TcChatType: db.In([]int{consts.ChatTypeMain, consts.ChatTypeOut}),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if deleteErr != nil {
			log.Errorf("[FsChatDisbandCallback] orgId: %d,  err: %v", orgId, deleteErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, deleteErr)
		}

		// 更新群聊设置
		_, updateErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, db.Cond{
			consts.TcChatId: chatId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsNoDelete,
		})
		if updateErr != nil {
			log.Errorf("[FsChatDisbandCallback] chatId: %s, err: %v", chatId, updateErr)
			return updateErr
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return errs.MysqlOperateError
	}

	if projectId != 0 {
		clearErr := domain.ClearProjectMainChatCache(orgId, projectId)
		if clearErr != nil {
			log.Error(clearErr)
			return clearErr
		}
	}

	return nil
}

// UpdateProjectDetailChatSetting 更新群聊关联，开启或关闭群聊
func UpdateProjectDetailChatSetting(orgId int64, projectId int64, isOpenChat int) errs.SystemErrorInfo {
	//查找项目对应的主群聊
	chatId, err := domain.GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if isOpenChat == 2 && chatId != "" {
			//删除所有与该群相关的项目关联
			_, deleteErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectChat, db.Cond{
				consts.TcOrgId:  orgId,
				consts.TcChatId: chatId,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
			})
			if deleteErr != nil {
				log.Errorf("[UpdateProjectDetailChatSetting] err: %v", deleteErr)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, deleteErr)
			}
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return errs.MysqlOperateError
	}

	if projectId != 0 {
		clearErr := domain.ClearProjectMainChatCache(orgId, projectId)
		if clearErr != nil {
			log.Error(clearErr)
			return clearErr
		}
	}

	return nil
}

func GetProjectMainChatId(orgId int64, currentUserId int64, projectId int64, sourceChannel string) (string, errs.SystemErrorInfo) {
	orgBaseInfo := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if orgBaseInfo.Failure() {
		log.Error(orgBaseInfo.Error())
		return "", orgBaseInfo.Error()
	}

	if orgBaseInfo.BaseOrgInfo.SourceChannel != sdk_const.SourceChannelFeishu {
		return "not open chat", nil
	}

	// 如果群聊是关闭的，则直接返回 not open chat
	//判断是否是项目成员
	isMember, isMemberErr := domain.IsProjectParticipant(orgId, currentUserId, projectId)
	if isMemberErr != nil {
		log.Error(isMemberErr)
		return "", isMemberErr
	}
	chatId, err := domain.GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Errorf("[GetProjectMainChatId] GetProjectMainChatId err: %v", err)
		return "", err
	}
	if chatId == "" {
		err := errs.HasNotProChat
		log.Errorf("[GetProjectMainChatId] err: %v, projectId: %d", err, projectId)
		return "", err
	}
	curUserResp := orgfacade.GetBaseUserInfo(orgvo.GetBaseUserInfoReqVo{
		OrgId:  orgId,
		UserId: currentUserId,
	})
	if curUserResp.Failure() {
		log.Errorf("[GetProjectMainChatId] GetBaseUserInfo err: %v, userId: %d", curUserResp.Error(), currentUserId)
		return "", curUserResp.Error()
	}
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[GetProjectMainChatId] GetBaseOrgInfoRelaxed err: %v", err)
		return "", err
	}
	projectInfo, err := domain.GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Errorf("[GetProjectMainChatId] GetProject err: %v", err)
		return "", err
	}

	// isSysAdmin := domain.CheckIsAdmin(orgId, currentUserId)
	isProAdmin, err := domain.CheckIsProAdmin(orgId, projectInfo.AppId, currentUserId, nil)
	if err != nil {
		log.Errorf("[GetProjectMainChatId] CheckIsProAdmin err: %v, appId: %d", err, projectInfo.AppId)
		return "", err
	}
	if !isProAdmin {
		// 检查用户是否在群聊中
		isInChat, err := domain.CheckUserIsInFsChat(baseOrgInfo, curUserResp.BaseUserInfo, chatId)
		if err != nil {
			log.Errorf("[GetProjectMainChatId] CheckUserIsInFsChat err: %v", err)
			return "", err
		}
		// 无论是项目成员，还是协作人，只要不在群中，都进行异常提示
		if !isInChat {
			return "", errs.ProChatNotInChat
		}
		if isMember {
			return chatId, nil
		}

		// 检查是否是协作人
		if !domain.CheckIsAppCollaborator(orgId, projectInfo.AppId, currentUserId) {
			return "", nil
		} else {
			return chatId, nil
		}
	} else {
		// 管理员则直接拉入群
		if err := domain.AddUserToChat(baseOrgInfo.OutOrgId, curUserResp.BaseUserInfo.OutUserId, chatId); err != nil {
			log.Errorf("[GetProjectMainChatId] AddUserToChat err: %v", err)
			return "", err
		}
	}

	return chatId, nil
}

// CheckIsShowProChatIcon 项目的任务列表页，检查是否是展示项目群聊 icon
// 对“项目成员”、“组织管理员”、“项目内的任务协作人”可见
func CheckIsShowProChatIcon(orgId int64, currentUserId int64, sourceChannel string, input projectvo.CheckIsShowProChatIconReqData) (bool,
	errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	// 目前只有飞书具备此功能
	if sourceChannel != sdk_const.SourceChannelFeishu {
		return false, nil
	}
	proAppId, oriErr := strconv.ParseInt(input.AppId, 10, 64)
	if oriErr != nil {
		log.Errorf("[CheckIsShowProChatIcon] err: %v, appId: %s", oriErr, input.AppId)
		return false, errs.BuildSystemErrorInfo(errs.TypeConvertError, oriErr)
	}
	projectInfo, err := domain.GetProjectByAppId(orgId, proAppId)
	if err != nil {
		log.Errorf("[CheckIsShowProChatIcon] GetProjectByAppId err: %v, proAppId: %d", err, proAppId)
		return false, err
	}
	// 检查是否是项目成员
	isMember, err := domain.IsProjectParticipant(orgId, currentUserId, projectInfo.Id)
	if err != nil {
		log.Errorf("[CheckIsShowProChatIcon] IsProjectParticipant err: %v", err)
		return false, err
	}
	if isMember {
		return true, nil
	}
	// 检查是否是系统管理员
	isSysAdmin := domain.CheckIsAdmin(orgId, currentUserId, proAppId)
	if isSysAdmin {
		return true, nil
	}
	// 检查是否是项目内的任务协作人
	if domain.CheckIsAppCollaborator(orgId, proAppId, currentUserId) {
		return true, err
	}

	return false, nil
}

func GetFsProjectChatPushSettings(orgId int64, input *projectvo.GetFsProjectChatPushSettingsReq) (*vo.GetFsProjectChatPushSettingsResp, errs.SystemErrorInfo,
) {
	res := vo.GetFsProjectChatPushSettingsResp{
		OutChatSettings: &vo.GetFsProjectChatPushSettingsOneChat{
			Tables:              make([]*vo.GetFsProjectChatPushSettingsOneChatTables, 0),
			ModifyColumnsOfSend: make([]*string, 0),
		},
	}
	if input.SourceChannel != sdk_const.SourceChannelFeishu {
		//此配置只对飞书用户有效
		return &res, errs.CannotBindChat
	}
	if input == nil {
		return &res, errs.ParamError
	}
	// 如果 chatId 为空，则通过 projectId 间接查询
	if input.ChatId == "" {
		// 查询 projectId 对应的主群聊
		chatId, err := domain.GetProjectMainChatId(orgId, input.ProjectId)
		if err != nil {
			log.Errorf("[ValidateAndFixInputForUpdateChatSetting] GetProjectMainChatId err: %v, proId: %d", err, input.ProjectId)
			return &res, err
		}
		input.ChatId = chatId
	}

	settings, err := domain.GetFsChatPushSettings(orgId, input.ChatId)
	if err != nil {
		log.Error(err)
		return &res, err
	}
	copyErr := copyer.Copy(settings, res.OutChatSettings)
	if copyErr != nil {
		log.Errorf("[GetFsProjectChatPushSettings] err: %v", copyErr)
		return &res, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &res, nil
}

// UpdateFsProjectChatPushSettings 更新某个群聊对应的推送配置
func UpdateFsProjectChatPushSettings(orgId int64, userId int64, sourceChannel string, params vo.UpdateFsProjectChatPushSettingsReq) errs.SystemErrorInfo {
	if sourceChannel != sdk_const.SourceChannelFeishu {
		return nil
	}
	if err := domain.ValidateAndFixInputForUpdateChatSetting(orgId, &params); err != nil {
		log.Errorf("[UpdateFsProjectChatPushSettings] ValidateAndFixInputForUpdateChatSetting err: %v", err)
		return err
	}
	// 保存表所属项目
	tableIdMapKeyByTableId := make(map[int64]int64, len(params.Tables))
	needGetTableAppIds, projectIds, tableIds := make([]int64, 0), make([]int64, 0), make([]int64, 0)
	for _, table := range params.Tables {
		projectIds = append(projectIds, table.ProjectID)
		tableId := int64(0)
		if table.TableID != nil {
			tableId = cast.ToInt64(*table.TableID)
		}
		if tableId > 0 {
			tableIds = append(tableIds, tableId)
			tableIdMapKeyByTableId[tableId] = table.ProjectID
		} else {
			// 如果客户端没传tableId，则证明需要所有table
			needGetTableAppIds = append(needGetTableAppIds, table.AppID)
		}
	}
	projectIds = slice.SliceUniqueInt64(projectIds)

	if len(needGetTableAppIds) > 0 {
		needGetTableAppIds = slice.SliceUniqueInt64(needGetTableAppIds)
		tablesMap, err := domain.GetTableListByProAppIds(orgId, needGetTableAppIds)
		if err != nil {
			log.Errorf("[UpdateFsProjectChatPushSettings] err: %v, chatId: %s", err, params.ChatID)
			return err
		}
		for _, data := range tablesMap {
			projectId := int64(0)
			for _, table := range params.Tables {
				if table.AppID == data.AppId {
					projectId = table.ProjectID
				}
			}
			if projectId > 0 {
				for _, datum := range data.Tables {
					tableIds = append(tableIds, datum.TableId)
					tableIdMapKeyByTableId[datum.TableId] = projectId
				}
			}
		}
	}

	// 查询现有的配置
	existSettingPoArr, err := domain.GetProChatSettingByChatId(orgId, params.ChatID)
	if err != nil {
		log.Errorf("[UpdateFsProjectChatPushSettings] err: %v, chatId: %s", err, params.ChatID)
		return err
	}
	settingMapByProId := make(map[int64]po.PpmProProjectChat, len(existSettingPoArr))

	// 查询主群聊关系
	proIdsForMain, err := domain.GetProjectIdsInChatIsMain(orgId, projectIds, params.ChatID)
	if err != nil {
		log.Errorf("[UpdateFsProjectChatPushSettings] err: %v, orgId: %d, chatID: %s", err, orgId, params.ChatID)
		return err
	}

	// 这批项目中，已存在的表群聊配置
	existTableIds := make([]int64, 0)
	for _, item := range existSettingPoArr {
		existTableIds = append(existTableIds, item.TableId)
		tableIdMapKeyByTableId[item.TableId] = item.ProjectId
		settingMapByProId[item.ProjectId] = item
	}
	_, add, del := int642.CompareSliceAddDelInt64(tableIds, existTableIds)
	needUpdateTableIds := int642.Int64Intersect(tableIds, existTableIds)
	err = domain.DealTransForUpdateGroupChatSettings(orgId, userId, &params, needUpdateTableIds, add, del, tableIdMapKeyByTableId, proIdsForMain)
	if err != nil {
		log.Errorf("[UpdateFsProjectChatPushSettings] DealTransForUpdateGroupChatSettings err: %v, orgId: %d, "+
			"chatID: %s", err, orgId, params.ChatID)
		return err
	}
	// 更新配置后，推送一条帮助卡片
	if err := domain.PushHelpCardAfterUpdateSetting(orgId, params.ChatID, projectIds); err != nil {
		log.Errorf("[UpdateFsProjectChatPushSettings] PushHelpCardAfterUpdateSetting err: %v, orgId: %d, "+
			"chatID: %s", err, orgId, params.ChatID)
		return err
	}

	return nil
}

// DeleteFsChat 删除主动创建的飞书群聊（创建项目时主动创建的群聊）
func DeleteFsChat(orgId int64, projectId int64, chatId string) errs.SystemErrorInfo {
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[DeleteFsChat] err: %v, projectId: %d", err, projectId)
		return err
	}
	if baseOrgInfo.OutOrgId == "" {
		return errs.CannotBindChat
	}
	tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err1 != nil {
		log.Errorf("[DeleteFsChat] err: %v, outOrgId: %s", err1, baseOrgInfo.OutOrgId)
		return err1
	}
	resp, respErr := tenant.DisbandChat(fsvo.UpdateChatData{ChatId: chatId})
	if respErr != nil {
		log.Error(respErr)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, respErr)
	}
	if resp.Code != 0 {
		log.Errorf("[DeleteFsChat] 删除飞书群聊失败 err: %s。outOrgId: %s", resp.Msg, baseOrgInfo.OutOrgId)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	// 删除群聊配置
	errSys := domain.DeleteGroupChatSettings(orgId, projectId, chatId)
	if errSys != nil {
		log.Errorf("[DeleteFsChat] DeleteGroupChatSettings failed, err: %v, chatId: %s", errSys, chatId)
		return errSys
	}
	// 清理缓存
	if err := domain.ClearProjectMainChatCache(orgId, projectId); err != nil {
		log.Errorf("[DeleteFsChat] ClearProjectMainChatCache err: %v, projectId: %d", err, projectId)
		return err
	}

	return nil
}

func AddDingChat(orgId, currentUserId int64, memberIds []int64, projectId int64, name string, ownerIds []int64, departmentIds []int64) (string, errs.SystemErrorInfo) {
	// 钉钉 获取成员的 userIds
	memberIds = append(memberIds, ownerIds...)
	//部门成员
	if len(departmentIds) > 0 {
		if ok, _ := slice.Contain(departmentIds, int64(0)); ok {
			userIdsResp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{OrgId: orgId})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return "", userIdsResp.Error()
			}
			memberIds = userIdsResp.Data
		} else {
			userIdsResp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: departmentIds,
			})
			if userIdsResp.Failure() {
				log.Error(userIdsResp.Error())
				return "", userIdsResp.Error()
			}
			memberIds = append(memberIds, userIdsResp.Data.UserIds...)
		}
	}
	memberIds = slice.SliceUniqueInt64(memberIds)

	userInfo := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		UserIds: memberIds,
		OrgId:   orgId,
	})
	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return "", userInfo.Error()
	}

	// 当前用户默认为群聊的群主，其他的项目成员是群里组员
	userIdsExcludeOwner := []string{}
	ownerUserId := ""
	for _, user := range userInfo.BaseUserInfos {
		if user.UserId == currentUserId {
			ownerUserId = user.OutUserId
			continue
		}
		userIdsExcludeOwner = append(userIdsExcludeOwner, user.OutUserId)
	}

	orgInfoBo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[AddDingChat]orgfacade.GetBaseOrgInfoRelaxed failed, orgId:%d, userId:%d ,projectId:%d, err:%v",
			orgId, currentUserId, projectId, err)
		return "", err
	}
	// 模版需要启用激活
	// 通过模版id创建钉钉的场景群
	client, errClient := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, orgInfoBo.OutOrgId)
	if errClient != nil {
		return "", errs.DingTalkClientError
	}
	ding := client.GetOriginClient().(*dingtalk.DingTalk)

	chat, errChat := ding.SceneGroupCreateChat(&request.SceneGroupCreatChat{
		Title:      name,
		TemplateId: "",
		Owner:      ownerUserId,
		Users:      strings.Join(userIdsExcludeOwner, ","),
	})
	if errChat != nil {
		return "", errs.DingTalkOpenApiCallError
	}

	openConversationId := chat.Result.OpenConversationId // 群id
	//chatId := chat.Result.ChatId                         // 会话id

	// 保存群id、以及会话id(客户端可能需要这个)
	//创建主群聊关系
	_, insertErr := domain.CreateMainChat(orgId, currentUserId, projectId, openConversationId, nil)
	if insertErr != nil {
		log.Error(insertErr)
		return "", insertErr
	}
	//插入表
	_, insertErr1 := domain.CreateChat(orgId, currentUserId, projectId, []string{openConversationId})
	if insertErr1 != nil {
		log.Error(insertErr1)
		return "", insertErr
	}

	// 发送群消息
	msg, errClient := ding.GroupSendSendAssistantMsg(&request.GroupSendSendAssistantMsgReq{
		TargetOpenConversationId: openConversationId,
		MsgTemplateId:            "", // 模版id
		MsgParamMap:              HandleDingReplyMsg(orgId, projectId),
		MsgMediaIdParamMap:       "",
		RobotCode:                "", // 机器人code
	})
	if errClient != nil {
		return "", errs.DingTalkOpenApiCallError
	}

	if msg.Code != 0 {
		log.Error(msg.ErrorMessage())
	}

	return openConversationId, nil
}

func DeleteChatCallback(outOrgId, chatId string, sourceChannel string) errs.SystemErrorInfo {
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: outOrgId})
	if orgInfoResp.Failure() {
		log.Error(orgInfoResp.Message)
		return orgInfoResp.Error()
	}
	orgId := orgInfoResp.BaseOrgInfo.OrgId

	// 删除群聊 chatId 对应的配置
	err := domain.DeleteGroupSettingsByChatId(orgId, chatId)
	if err != nil {
		log.Errorf("[DeleteChatCallback]failed, err:%v, orgId:%v, chatId:%s, sourceChannel:%s", err, orgId, chatId, sourceChannel)
		return err
	}

	return nil
}

// 消息模版按钮的链接替换
func HandleDingReplyMsg(orgId int64, projectId int64) string {
	//dingConfig := config.GetConfig().DingTalk

	return ""
}
