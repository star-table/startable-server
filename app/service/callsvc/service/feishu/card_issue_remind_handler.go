package callsvc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/app/facade/pushfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/types"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

type IssueRemindHandler struct{}

const SourceChannel = sdk_const.SourceChannelFeishu

// 仅仅在飞书平台，即将逾期提醒的卡片回调，交互字段：任务状态和截止时间
func (IssueRemindHandler) Handle(cardReq CardReq) (string, errs.SystemErrorInfo) {
	log.Infof("[IssueRemindHandler] handle info:%s", json.ToJsonIgnoreError(cardReq))
	value := cardReq.Action.Value

	// 操作人的 openId
	opUserOpenId := cardReq.OpenId
	//tenantKey := cardReq.TenantKey
	optionsDate := cardReq.Action.Option
	issueId := value.IssueId
	statusType := value.StatusType
	action := value.Action

	orgId := value.OrgId
	tableId, err := strconv.ParseInt(value.TableId, 10, 64)
	if err != nil {
		log.Errorf("[IssueRemindHandler][Handle] parse tableId err: %v, tableIdStr: %s", err, value.TableId)
		return "", errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
	}
	var oriErr error
	appId := int64(0)
	appId, oriErr = strconv.ParseInt(value.AppId, 10, 64)
	if oriErr != nil {
		log.Errorf("[IssueRemindHandler] parse appId err:%v, orgId:%d, appId: %s", oriErr, orgId, value.AppId)
		return "", errs.TypeConvertError
	}

	log.Infof("[IssueRemindHandler][Handle] 卡片通知，编辑任务回调，任务id %d, 更改的截止时间 %s, 更改的状态类型 %d", issueId, optionsDate, statusType)
	// 查询操作人信息
	resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: opUserOpenId,
	})
	if resp.Failure() {
		return "", errs.BuildSystemErrorInfo(errs.UserNotFoundError, resp.Error())
	}
	opUserInfo := resp.BaseUserInfo
	//userId := value.UserId
	userId := opUserInfo.UserId //操作人就是传入的操作人

	issueInfoList := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
		UserId:   userId,
		OrgId:    orgId,
		IssueIds: []int64{issueId},
	})
	if issueInfoList.Failure() {
		log.Errorf("[IssueRemindHandler]handler GetIssueInfoList failed:%v, orgId:%d, userId:%d",
			issueInfoList.Error(), orgId, userId)
		return "", issueInfoList.Error()
	}
	issueInfo := issueInfoList.IssueInfos[0]
	projectId := issueInfo.ProjectID

	if orgId == 0 || userId == 0 {
		log.Errorf("卡片通知，编辑任务回调，任务id %d, 用户或组织信息未指定", issueId)
		return "", errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)
	}

	// 项目/表
	projectContent := ""
	projectName := consts.CardDefaultIssueProjectName
	//projectTypeId := int64(0)
	if projectId > 0 {
		projectInfo := projectfacade.ProjectInfo(projectvo.ProjectInfoReqVo{
			Input:         vo.ProjectInfoReq{ProjectID: projectId},
			OrgId:         orgId,
			UserId:        userId,
			SourceChannel: sdk_const.SourceChannelFeishu,
		})
		if projectInfo.Failure() {
			log.Errorf("[IssueRemindHandler]handler projectInfo failed:%v, orgId:%d, userId:%d",
				projectInfo.Error(), orgId, userId)
			return "", projectInfo.Error()
		}
		projectName = projectInfo.ProjectInfo.Name
		//projectTypeId = projectInfo.ProjectInfo.ProjectTypeID
		appIdStr := projectInfo.ProjectInfo.AppID
		parseInt, err := strconv.ParseInt(appIdStr, 10, 64)
		if err != nil {
			log.Errorf("[IssueRemindHandler]handler failed:%v, orgId:%d, userId:%d", err, orgId, userId)
			return "", errs.TypeConvertError
		}
		appId = parseInt
	} else {
		var oriErr error
		appId, oriErr = strconv.ParseInt(value.AppId, 10, 64)
		if oriErr != nil {
			log.Errorf("[IssueRemindHandler] parse appId err:%v, orgId:%d, userId:%d, appId: %s", oriErr,
				orgId, userId, value.AppId)
			return "", errs.TypeConvertError
		}
		// 查汇总表
		orgResp := orgfacade.GetOrgBoListByPage(orgvo.GetOrgIdListByPageReqVo{
			Page:  1,
			Size:  1,
			Input: orgvo.GetOrgIdListByPageReqVoData{OrgIds: []int64{orgId}},
		})
		if orgResp.Failure() {
			log.Errorf("[IssueRemindHandler] GetOrgBoListByPage err:%v, orgId:%d, userId:%d, appId: %v",
				orgResp.Error(), orgId, userId, value.AppId)
			return "", orgResp.Error()
		}
		if len(orgResp.Data.List) < 1 {
			return "", nil
		}
		org := orgResp.Data.List[0]
		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
		if len(org.Remark) > 0 {
			oriErr = json.FromJson(org.Remark, orgRemarkObj)
			if oriErr != nil {
				log.Errorf("[IssueRemindHandler] 组织 remark 反序列化异常，组织id:%d,原因:%v", org.Id, oriErr)
				return "", errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
			}
		}
		appId = orgRemarkObj.OrgSummaryTableAppId
	}

	//tableIdInt64 := int64(0)
	tableName := ""
	if projectName != consts.CardDefaultIssueProjectName {
		tableIdStr := issueInfo.TableID
		tableId, err := strconv.ParseInt(tableIdStr, 10, 64)
		if err != nil {
			log.Errorf("[IssueRemindHandler]handler failed:%v, orgId:%d, userId:%d", err, orgId, userId)
			return "", errs.TypeConvertError
		}
		//tableIdInt64 = tableId

		respTable := projectfacade.GetTable(projectvo.GetTableInfoReq{
			OrgId:  orgId,
			UserId: userId,
			Input:  &tablePb.ReadTableRequest{TableId: tableId},
		})
		if respTable.Failure() {
			log.Errorf("[IssueRemindHandler]handler GetTable failed:%v, orgId:%d, userId:%d",
				respTable.Error(), orgId, userId)
		}

		tableName = respTable.Data.Table.Name

		projectContent = fmt.Sprintf(consts.CardTablePro, projectName, tableName)
	} else {
		projectContent = consts.CardDefaultIssueProjectName
	}

	// 父记录
	parentTitle := consts.CardDefaultRelationIssueTitle
	if issueInfo.ParentID > 0 {
		parents := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
			UserId:   userId,
			OrgId:    orgId,
			IssueIds: []int64{issueInfo.ParentID},
		})
		if parents.Failure() {
			log.Errorf("[IssueRemindHandler]handler error:%v, orgId:%d, parentId:%d",
				parents.Error(), orgId, issueInfo.ParentID)
			return "", parents.Error()
		}
		if len(parents.IssueInfos) > 0 {
			parentTitle = parents.IssueInfos[0].Title
		}
	}

	// 负责人
	owners := []string{}
	ownerInfos := []*bo.BaseUserInfoBo{}
	if len(issueInfo.Owners) > 0 {
		batchUsers := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
			OrgId:   orgId,
			UserIds: issueInfo.Owners,
		})
		if batchUsers.Failure() {
			log.Errorf("[IssueRemindHandler]handler GetBaseUserInfoBatch failed:%v, orgId:%d, userId:%d",
				batchUsers.Failure(), orgId, userId)
			return "", batchUsers.Error()
		}
		for _, owner := range batchUsers.BaseUserInfos {
			owners = append(owners, owner.Name)
			ownerInfos = append(ownerInfos, &owner)
		}
	}

	ownerDisplayName := strings.Join(owners, "，")
	if ownerDisplayName == "" {
		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	}

	issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         orgId,
		IssueId:       issueId,
	}).Data

	columnsResp := projectfacade.GetOneTableColumns(projectvo.GetTableColumnReq{
		OrgId:     orgId,
		UserId:    userId,
		ProjectId: projectId,
		TableId:   tableId,
	})
	if columnsResp.Failure() {
		log.Errorf("[PushCreateIssueSuccessNotice] err:%v", columnsResp.Error())
		return "", columnsResp.Error()
	}

	columnMap := make(map[string]*projectvo.TableColumnData)
	for _, column := range columnsResp.Data.Columns {
		columnMap[column.Name] = column
	}

	if (action == consts.FsCardActionUpdatePlanEndTime || action == consts.FsCardActionUpdatePlanStartTime) && optionsDate != "" {
		targetTime, parseErr := time.Parse(consts.AppTimeFormatYYYYMMDDHHmmTimezone, optionsDate)
		if parseErr != nil {
			log.Error(parseErr)
			return "", errs.BuildSystemErrorInfo(errs.DateParseError)
		}
		targetTypesTime := types.Time(targetTime)
		form := make([]map[string]interface{}, 0)
		data := make(map[string]interface{})
		data[consts.BasicFieldId] = issueId
		//updateInput := vo.UpdateIssueReq{
		//	ID: issueId,
		//}
		if action == consts.FsCardActionUpdatePlanEndTime {
			data[consts.BasicFieldPlanEndTime] = targetTypesTime
			//updateInput.PlanEndTime = &targetTypesTime
			//updateInput.UpdateFields = []string{"planEndTime"}
		} else {
			data[consts.BasicFieldPlanStartTime] = targetTypesTime
			//updateInput.PlanStartTime = &targetTypesTime
			//updateInput.UpdateFields = []string{"planStartTime"}
		}
		form = append(form, data)
		updateRespVo := projectfacade.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqVo{
			UserId: opUserInfo.UserId,
			OrgId:  orgId,
			Input: &projectvo.BatchUpdateIssueInput{
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				Data:      form,
			},
		})
		if updateRespVo.Failure() {
			//if updateRespVo.Code == errs.NoOperationPermissionForIssueUpdate.Code() {
			// 群聊 && 更新任务时提示异常，则返回异常提示信息卡片。
			planEndTime := time.Time(issueInfo.PlanEndTime).Format(consts.AppTimeFormatYYYYMMDDHHmm)
			respCard := card.GetCardBeOverdue(&projectvo.CardBeOverdue{
				OrgId:          orgId,
				OwnerId:        userId,
				IssueId:        issueId,
				IssueLinks:     issueLinks,
				IssueTitle:     issueInfo.Title,
				PlanEndTime:    planEndTime,
				SourceChannel:  sdk_const.SourceChannelFeishu,
				TableId:        value.TableId,
				AppId:          value.AppId,
				ContentProject: projectContent,
				OwnerNameStr:   ownerDisplayName,
				ParentTitle:    parentTitle,
				HasPermission:  false,
				Tips:           updateRespVo.Message,
				TableColumn:    columnMap,
			})

			cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
				SourceChannel: sdk_const.SourceChannelFeishu,
				Card:          respCard,
			})
			if cardResp.Failure() {
				log.Errorf("[IssueRemindHandler]handler pushfacade.GenerateCard err:%v", cardResp.Error())
				return "", cardResp.Error()
			}
			return cardResp.Data.Card, nil

		}

		log.Infof("编辑任务回调，更新截止时间成功，任务id %d, 更改的截止时间 %s, 更改的状态类型 %d", issueId, optionsDate, statusType)
		lastTime := targetTime.Format(consts.AppTimeFormat)
		respCard := card.GetCardBeOverdue(&projectvo.CardBeOverdue{
			OrgId:          orgId,
			OwnerId:        userId,
			IssueId:        issueId,
			IssueLinks:     issueLinks,
			IssueTitle:     issueInfo.Title,
			PlanEndTime:    lastTime,
			SourceChannel:  sdk_const.SourceChannelFeishu,
			TableId:        value.TableId,
			AppId:          value.AppId,
			ContentProject: projectContent,
			OwnerNameStr:   ownerDisplayName,
			ParentTitle:    parentTitle,
			HasPermission:  true,
			Tips:           "",
			TableColumn:    columnMap,
		})
		cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
			SourceChannel: sdk_const.SourceChannelFeishu,
			Card:          respCard,
		})
		if cardResp.Failure() {
			log.Errorf("[IssueRemindHandler]handler pushfacade.GenerateCard err:%v", cardResp.Error())
			return "", cardResp.Error()
		}
		return cardResp.Data.Card, nil
	}
	if action == consts.FsCardActionUpdateIssueStatus && statusType > 0 && statusType < 4 {
		form := make([]map[string]interface{}, 0)
		data := make(map[string]interface{})
		data[consts.BasicFieldId] = issueId
		data[consts.BasicFieldIssueStatusType] = statusType
		form = append(form, data)

		updateResp := projectfacade.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqVo{
			UserId: opUserInfo.UserId,
			OrgId:  orgId,
			Input: &projectvo.BatchUpdateIssueInput{
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				Data:      form,
			},
		})

		if updateResp.Failure() {
			log.Errorf("编辑任务回调，更新状态类型失败，任务id %d, 更改的截止时间 %s, 更改的状态类型 %d,err:%v", issueId, optionsDate, statusType, updateResp.Error())

			//if updateResp.Code == errs.NoOperationPermissionForIssueUpdate.Code() {
			// 返回没有相关权限的卡片提示
			planEndTime := time.Time(issueInfo.PlanEndTime).Format(consts.AppTimeFormatYYYYMMDDHHmm)
			respCard := card.GetCardBeOverdue(&projectvo.CardBeOverdue{
				OrgId:          orgId,
				OwnerId:        userId,
				IssueId:        issueId,
				IssueLinks:     issueLinks,
				IssueTitle:     issueInfo.Title,
				PlanEndTime:    planEndTime,
				SourceChannel:  sdk_const.SourceChannelFeishu,
				TableId:        value.TableId,
				AppId:          value.AppId,
				ContentProject: projectContent,
				OwnerNameStr:   ownerDisplayName,
				ParentTitle:    parentTitle,
				HasPermission:  false,
				Tips:           updateResp.Message,
				TableColumn:    columnMap,
			})
			cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
				SourceChannel: sdk_const.SourceChannelFeishu,
				Card:          respCard,
			})
			if cardResp.Failure() {
				log.Errorf("[IssueRemindHandler]handler pushfacade.GenerateCard err:%v", cardResp.Error())
				return "", cardResp.Error()
			}
			log.Infof("[IssueRemindHandler]handler FsCardBeOverdueInfoWithoutPermission: %s", cardResp.Data.Card)
			return cardResp.Data.Card, nil
			//}
			//return "", updateResp.Error()
		}

		log.Infof("[IssueRemindHandler]handler 编辑任务回调，更新状态类型成功，任务id %d, 更改的截止时间 %s, 更改的状态类型 %d", issueId,
			optionsDate, statusType)

		respCard := card.GetFsCardBeOverdueComplete(issueInfo.Title, projectId, projectName, tableName, ownerInfos, issueLinks, columnMap)
		cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
			SourceChannel: sdk_const.SourceChannelFeishu,
			Card:          respCard,
		})
		if cardResp.Failure() {
			log.Errorf("[IssueRemindHandler]handler pushfacade.GenerateCard err:%v", cardResp.Error())
			return "", cardResp.Error()
		}
		return cardResp.Data.Card, nil
	}

	return "", nil
}

//// 群聊中，发送卡片，提示用户信息
//func ReplyTipCard(tenantKey, tipText, openChatId, targetOpenId, targetUserName string) error {
//	tenant, err := feishu.GetTenant(tenantKey)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	msgCard := GetCardForTip(tipText, map[string]interface{}{
//		"atOpenId":   targetOpenId,
//		"atUserName": targetUserName,
//	})
//	fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
//		ChatId:  openChatId,
//		MsgType: "interactive",
//		Card:    msgCard,
//	})
//	if oriErr != nil {
//		log.Error(oriErr)
//		return err
//	}
//	if fsResp.Code != 0 {
//		log.Error("HandleGroupChatAtUserName isIn 发送消息异常")
//	}
//	return nil
//}

///// 回复提示信息卡片，如：抱歉，您指定的用户不是项目成员。
//func GetCardForTip(tipText string, extraParam map[string]interface{}) *fsvo.Card {
//	// 实现 at 某个成员
//	openId := ""
//	userName := ""
//	if tmpVal, ok := extraParam["atOpenId"]; ok {
//		openId = tmpVal.(string)
//	}
//	if tmpVal, ok := extraParam["atUserName"]; ok {
//		userName = tmpVal.(string)
//	}
//	atSomeoneStr := str.RenderAtSomeoneStr(openId, userName)
//	if len(openId) > 0 && len(atSomeoneStr) > 0 {
//		tipText = atSomeoneStr + tipText
//	}
//	cardMsg := &fsvo.Card{
//		Header: &fsvo.CardHeader{
//			Title: &fsvo.CardHeaderTitle{
//				Tag:     "plain_text",
//				Content: "提示信息",
//			},
//		},
//		Elements: []interface{}{
//			fsvo.CardElementContentModule{
//				Tag: "div",
//				Fields: []fsvo.CardElementField{
//					{
//						Text: fsvo.CardElementText{
//							Tag:     "lark_md",
//							Content: tipText,
//						},
//					},
//				},
//			},
//		},
//	}
//
//	return cardMsg
//}
