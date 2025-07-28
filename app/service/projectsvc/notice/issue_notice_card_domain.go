package notice

import (
	"fmt"
	"strings"

	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"google.golang.org/protobuf/types/known/structpb"
)

// GetIssueCardNormal 组装飞书卡片：创建任务、更新任务、删除任务等
func GetIssueCardNormal(issueTrendsBo bo.IssueTrendsBo) (*pushPb.TemplateCard, error) {
	orgId := issueTrendsBo.OrgId
	operatorId := issueTrendsBo.OperatorId
	//pushType := issueNoticeBo.PushType

	ownerIds := issueTrendsBo.AfterOwner
	userIds := []int64{}
	userIds = append(userIds, ownerIds...)
	userIds = append(userIds, operatorId)
	userIds = slice.SliceUniqueInt64(userIds)

	baseUserInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Errorf("GetIssueCardNormal GetBaseUserInfoBatch err:%v, orgId:%v, userIds:%v", err, orgId, userIds)
		return nil, err
	}

	userMap := make(map[int64]bo.BaseUserInfoBo, len(baseUserInfos))
	for _, u := range baseUserInfos {
		userMap[u.UserId] = u
	}

	ownerInfos := make([]*bo.BaseUserInfoBo, 0, len(baseUserInfos))
	operatorBaseInfo := userMap[operatorId]

	for _, uId := range issueTrendsBo.AfterOwner {
		ownerInfo := userMap[uId]
		ownerInfos = append(ownerInfos, &ownerInfo)
	}

	// 负责人
	ownerSlice := []string{}
	for _, ownerInfo := range ownerInfos {
		ownerSlice = append(ownerSlice, ownerInfo.Name)
	}
	ownerDisplayName := strings.Join(ownerSlice, "，")
	if ownerDisplayName == "" {
		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	}

	projectInfo, err1 := domain.LoadProjectAuthBo(orgId, issueTrendsBo.ProjectId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	issueInfo, err := domain.GetIssueInfoLc(orgId, 0, issueTrendsBo.IssueId)
	if err != nil {
		log.Errorf("[GetIssueCardNormal] domain.GetIssueInfoLc err:%v, orgId:%v, issueId:%v", err, orgId, issueTrendsBo.IssueId)
		return nil, err
	}
	//任务标题为空，则不需要推送
	if issueInfo.Title == "" {
		return nil, errs.CardTitleEmpty
	}
	taskTypeName := "记录"

	title := card.GetOperationName(&issueTrendsBo, operatorBaseInfo.Name) + " 创建了新的" + taskTypeName

	meta := &pushPb.TemplateCard{}
	meta.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: title,
	}

	links := domain.GetIssueLinks("", orgId, issueTrendsBo.IssueId)
	issueInfoAppLink := links.SideBarLink
	issueLink := links.Link

	tableColumns, err := domain.GetTableColumnsMap(orgId, issueTrendsBo.TableId, []string{consts.BasicFieldTitle, consts.BasicFieldOwnerId})
	if err != nil {
		log.Errorf("[GetIssueCardNormal] GetTableColumnsMap failed:%v, userId:%d, tableId:%d", err, issueTrendsBo.OperatorId, issueTrendsBo.TableId)
		return nil, err
	}

	columnTitle, columnOwnerId := card.GetCardTitleAndOwnerColumn(tableColumns)
	tableName := consts.DefaultTableName
	if issueTrendsBo.TableId != 0 {
		tableInfo, err := domain.GetTableByTableId(orgId, issueTrendsBo.OperatorId, issueTrendsBo.TableId)
		if err != nil {
			log.Errorf("[GetIssueCardNormal] GetTableInfo failed:%v, userId:%d, tableId:%d", err, issueTrendsBo.OperatorId, issueTrendsBo.TableId)
			return nil, err
		}
		tableName = tableInfo.Name
	}

	meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnTitle,
				Value: issueTrendsBo.IssueTitle,
			},
		},
	})

	meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementProjectTable,
				Value: fmt.Sprintf(consts.CardTablePro, projectInfo.Name, tableName),
			},
		},
	})

	if issueTrendsBo.ParentId != 0 {
		parent, err := domain.GetIssueInfoLc(orgId, 0, issueTrendsBo.ParentId)
		if err != nil {
			log.Errorf("[GetIssueCardNormal] GetIssueBo err: %v, parentId: %d", err, issueTrendsBo.ParentId)
			return nil, err
		}
		meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementParent,
					Value: parent.Title,
				},
			},
		})
	}

	meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnOwnerId,
				Value: ownerDisplayName,
			},
		},
	})

	m := map[string]interface{}{
		consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
		consts.FsCardValueAction:   consts.FsCardUnsubscribeIssueRemind,
		consts.FsCardValueOrgId:    orgId,
	}
	values, errSys := structpb.NewStruct(m)
	if errSys != nil {
		log.Errorf("[GetIssueCardNormal]err:%v", errSys)
		return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, errSys)
	}
	meta.ActionModules = []*pushPb.CardActionModule{
		&pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_Default,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewDetail,
					Url:   issueInfoAppLink,
				},
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_NotImportant,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewInsideApp,
					Url:   issueLink,
				},
			},
		},
		&pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Hr,
		},
		&pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level:  pushPb.CardElementLevel_Default,
					Type:   pushPb.CardActionType_Link,
					Text:   consts.CardUnSubscribe,
					Values: values,
				},
			},
		},
	}

	return meta, nil
}

// GetIssueCardUpdateIssue 组装更新任务时的卡片
func GetIssueCardUpdateIssue(issueTrendsBo *bo.IssueTrendsBo, infoObj *projectvo.BaseInfoBoxForIssueCard, cardMeta *pushPb.TemplateCard) (*pushPb.TemplateCard, *pushPb.TemplateCard, errs.SystemErrorInfo) {
	//任务标题为空，则不需要推送
	if infoObj.IssueInfo.Title == "" {
		return nil, nil, errs.CardTitleEmpty
	}

	taskTypeName := "记录"
	title := card.GetOperationName(issueTrendsBo, infoObj.OperateUser.Name) + " 更新了" + taskTypeName

	log.Infof("[GetIssueCardUpdateIssue] cardElementsForUpdateDetail len: %d", len(cardMeta.Divs))

	meta := card.GenerateCard(&commonvo.GenerateCard{
		OrgId:         issueTrendsBo.OrgId,
		CardTitle:     title,
		CardLevel:     pushPb.CardElementLevel_Default,
		IssueTitle:    issueTrendsBo.IssueTitle,
		ParentId:      issueTrendsBo.ParentId,
		ProjectId:     infoObj.ProjectAuthBo.Id,
		ProjectTypeId: infoObj.ProjectAuthBo.ProjectTypeId,
		ParentTitle:   infoObj.ParentIssue.Title,
		ProjectName:   infoObj.ProjectAuthBo.Name,
		TableName:     infoObj.IssueTableInfo.Name,
		OwnerInfos:    infoObj.OwnerInfos,
		Links: &projectvo.IssueLinks{
			Link:        infoObj.IssuePcUrl,
			SideBarLink: infoObj.IssueInfoUrl,
		},
		IsNeedUnsubscribe: true,
		SpecialDiv:        cardMeta.Divs, // 如果有更新详情，则追加上更新详情
		TableColumnInfo:   infoObj.ProjectTableColumn,
	})

	return meta, meta, nil
}

// GetIssueCardOwnerChange 任务负责人变更时的推送
// 返回两种卡片，针对当事人和非当事人
// memberType 成员类型：auditor：确认人；follower：关注人；owner: 负责人；
func GetIssueCardOwnerChange(issueTrendsBo *bo.IssueTrendsBo, issueNoticeBo bo.IssueNoticeBo, memberType string) (*pushPb.TemplateCard, *pushPb.TemplateCard, errs.SystemErrorInfo) {
	orgId := issueNoticeBo.OrgId
	ownerIds := issueNoticeBo.AfterOwner
	operatorId := issueNoticeBo.OperatorId
	projectId := issueNoticeBo.ProjectId

	userIds := make([]int64, 0, len(ownerIds)+1)
	userIds = append(userIds, operatorId)
	userIds = append(userIds, ownerIds...)
	usersMap := map[int64]bo.BaseUserInfoBo{}
	userInfoArr, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Errorf("查询组织 %d 用户 %d 信息出现异常 %v", orgId, operatorId, err)
		return nil, nil, err
	}

	for _, u := range userInfoArr {
		usersMap[u.UserId] = u
	}

	projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Errorf("[GetIssueCardOwnerChange] err: %v, projectId: %d", err, projectId)
		return nil, nil, err
	}

	links := domain.GetIssueLinks("", orgId, issueNoticeBo.IssueId)
	issueInfoAppLink := links.SideBarLink
	issueLink := links.Link
	tableName := consts.DefaultTableName
	if issueNoticeBo.TableId != 0 {
		tableInfo, err := domain.GetTableByTableId(orgId, issueNoticeBo.OperatorId, issueNoticeBo.TableId)
		if err != nil {
			log.Errorf("[GetIssueCardOwnerChange] GetTableInfo failed:%v, userId:%d, tableId:%d", err, issueNoticeBo.OperatorId, issueNoticeBo.TableId)
			return nil, nil, err
		}
		tableName = tableInfo.Name
	}

	memberDesc := "负责人"
	switch memberType {
	case "follower":
		memberDesc = "关注人"
	case "auditor":
		memberDesc = "确认人"
	}

	meta := &pushPb.TemplateCard{}

	// 标题
	meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementTitle,
				Value: issueNoticeBo.IssueTitle,
			},
		},
	})

	// 项目/表
	contentStr := ""
	if projectId == 0 {
		contentStr = consts.CardDefaultIssueProjectName
	} else {
		contentStr = fmt.Sprintf(consts.CardTablePro, projectInfo.Name, tableName)
	}
	meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementProjectTable,
				Value: contentStr,
			},
		},
	})

	if issueNoticeBo.ParentId != 0 {
		parents, err := domain.GetIssueInfosLc(orgId, 0, []int64{issueNoticeBo.ParentId})
		if err != nil {
			log.Errorf("[GetIssueCardOwnerChange] err: %v, issueNoticeBo.ParentId: %d", err, issueNoticeBo.ParentId)
			return nil, nil, err
		}
		parent := parents[0]
		parentTitle := parent.Title
		if parent.Title == "" {
			parentTitle = consts.CardDefaultRelationIssueTitle
		}
		meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementParent,
					Value: parentTitle,
				},
			},
		})
	}

	// 负责人
	ownerSlice := []string{}
	ownerInfos := []bo.BaseUserInfoBo{}
	for _, ownerId := range ownerIds {
		ownerInfos = append(ownerInfos, usersMap[ownerId])
	}
	for _, ownerInfo := range ownerInfos {
		ownerSlice = append(ownerSlice, ownerInfo.Name)
	}
	ownerDisplayName := strings.Join(ownerSlice, "，")
	if ownerDisplayName == "" {
		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	}
	meta.Divs = append(meta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementOwner,
				Value: ownerDisplayName,
			},
		},
	})
	// 操作人
	operatorName := usersMap[operatorId].Name

	// 负责人需展示起止时间 todo
	if memberType == "owner" {
	}

	// 卡片的按钮和退订按钮
	m := map[string]interface{}{
		consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
		consts.FsCardValueAction:   consts.FsCardUnsubscribeIssueRemind,
		consts.FsCardValueOrgId:    orgId,
	}
	values, errSys := structpb.NewStruct(m)
	if errSys != nil {
		log.Errorf("[GetIssueCardOwnerChange] err:%v", errSys)
		return nil, nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
	}
	meta.ActionModules = []*pushPb.CardActionModule{
		&pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_Default,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewDetail,
					Url:   issueInfoAppLink,
				},
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_NotImportant,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewInsideApp,
					Url:   issueLink,
				},
			},
		},
		&pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Hr,
		},
		&pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level:  pushPb.CardElementLevel_Default,
					Type:   pushPb.CardActionType_Link,
					Text:   consts.CardUnSubscribe,
					Values: values,
				},
			},
		},
	}

	meta.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: fmt.Sprintf("%s 将你添加为 %s", card.GetOperationName(issueTrendsBo, operatorName), memberDesc),
	}

	return meta, meta, nil
}

func GetIssueCardRemarkAt(issueNoticeBo bo.IssueNoticeBo, content string) (*pushPb.TemplateCard, *pushPb.TemplateCard, error) {
	orgId := issueNoticeBo.OrgId
	comment := GetFeiShuAtConent(orgId, content)
	table, err := domain.GetTableByTableId(orgId, issueNoticeBo.OperatorId, issueNoticeBo.TableId)
	if err != nil {
		log.Errorf("[GetIssueCardRemarkAt] err: %v", err)
		return nil, nil, err
	}

	links := domain.GetIssueLinks("", orgId, issueNoticeBo.IssueId)
	issueInfoAppLink := links.SideBarLink
	issueLink := links.Link

	projectInfo, err := domain.GetProjectSimple(orgId, issueNoticeBo.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(orgId, issueNoticeBo.OperatorId)
	if err != nil {
		log.Errorf("查询组织 %d 用户 %d 信息出现异常 %v", orgId, issueNoticeBo.OperatorId, err)
		return nil, nil, err
	}
	titleForSelf := "@我"
	titleForOther := fmt.Sprintf("%s 更新了记录描述", operatorBaseInfo.Name)

	markdownProject := fmt.Sprintf("%s/%s", projectInfo.Name, table.Name)

	commDivs := []*pushPb.CardTextDiv{}
	commDivs = append(commDivs,
		&pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementTitle,
					Value: issueNoticeBo.IssueTitle,
				},
			},
		},
		&pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementProjectTable,
					Value: markdownProject,
				},
			},
		},
		&pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   "操作人",
					Value: operatorBaseInfo.Name,
				},
			},
		},
		&pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   operatorBaseInfo.Name,
					Value: comment,
				},
			},
		},
	)

	commActions := []*pushPb.CardActionModule{}

	commFsAction := []interface{}{}
	commFsAction = append(commFsAction,
		vo.CardElementActionModule{
			Tag:    "action",
			Layout: "bisected",
			Actions: []interface{}{
				vo.ActionButton{
					Tag: "button",
					Text: vo.CardElementText{
						Tag:     "plain_text",
						Content: consts.CardButtonTextForViewDetail,
					},
					Url:  issueInfoAppLink,
					Type: consts.FsCardButtonColorPrimary,
				},
				vo.ActionButton{
					Tag: "button",
					Text: vo.CardElementText{
						Tag:     "plain_text",
						Content: consts.CardButtonTextForViewInsideApp,
					},
					Url:  issueLink,
					Type: consts.FsCardButtonColorDefault,
				},
			},
		},
		vo.CardElementActionModule{
			Tag: "hr",
		},
		vo.CardElementActionModule{
			Tag: "action",
			Actions: []interface{}{
				vo.ActionButton{
					Tag: "button",
					Text: vo.CardElementText{
						Tag:     "plain_text",
						Content: "点击退订",
					},
					Type: "link",
					Value: map[string]interface{}{
						consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
						consts.FsCardValueAction:   consts.FsCardUnsubscribeIssueRemind,
						consts.FsCardValueOrgId:    orgId,
					},
				},
			},
		},
	)

	cardForSelf := &pushPb.TemplateCard{
		Title: &pushPb.CardTitle{
			Level: pushPb.CardElementLevel_Default,
			Title: titleForSelf,
		},
		Divs:          commDivs,
		ActionModules: commActions,
	}

	cardForOther := &pushPb.TemplateCard{
		Title: &pushPb.CardTitle{
			Level: pushPb.CardElementLevel_Default,
			Title: titleForOther,
		},
		Divs:          commDivs,
		ActionModules: commActions,
	}

	return cardForSelf, cardForOther, nil
}

// GetIssueCardComment 组装评论时的卡片
func GetIssueCardComment(issueNoticeBo bo.IssueNoticeBo, infoBox *projectvo.BaseInfoBoxForIssueCard, content string) (*pushPb.TemplateCard, *pushPb.TemplateCard, error) {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(orgId, operatorId)
	if err != nil {
		log.Errorf("[GetIssueCardComment] 查询组织 %d 用户 %d 信息出现异常 %v", orgId, operatorId, err)
		return nil, nil, err
	}
	links := domain.GetIssueLinks("", orgId, issueNoticeBo.IssueId)
	comment := GetFeiShuAtConent(orgId, content)

	titleForSelf := fmt.Sprintf("%s@我", operatorBaseInfo.Name)
	titleForOther := fmt.Sprintf("%s 评论了记录", operatorBaseInfo.Name)
	comment = domain.CutComment(comment, consts.GroupChatIssueCommentLimit)
	speakerName := fmt.Sprintf("%s:", operatorBaseInfo.Name)

	markdownComment := fmt.Sprintf("%s", "“"+comment+"”")

	cardDivs := []*pushPb.CardTextDiv{
		&pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   speakerName,
					Value: markdownComment,
				},
			},
		},
	}
	cardCommon := card.GenerateCard(&commonvo.GenerateCard{
		OrgId: orgId,
		//CardTitle:         titleForSelf,
		CardLevel:         pushPb.CardElementLevel_Default,
		ProjectTypeId:     infoBox.ProjectAuthBo.ProjectTypeId,
		ProjectName:       infoBox.ProjectAuthBo.Name,
		TableName:         infoBox.IssueTableInfo.Name,
		IssueTitle:        infoBox.IssueInfo.Title,
		ParentId:          infoBox.IssueInfo.ParentId,
		ProjectId:         infoBox.IssueInfo.ProjectId,
		OwnerInfos:        infoBox.OwnerInfos,
		Links:             links,
		IsNeedUnsubscribe: true,
		SpecialDiv:        cardDivs,
		TableColumnInfo:   infoBox.ProjectTableColumn,
	})

	cardForSelf := &pushPb.TemplateCard{
		Title: &pushPb.CardTitle{
			Level: pushPb.CardElementLevel_Default,
			Title: titleForSelf,
		},
		Divs:          cardCommon.Divs,
		ActionModules: cardCommon.ActionModules,
	}

	cardForOther := &pushPb.TemplateCard{
		Title: &pushPb.CardTitle{
			Level: pushPb.CardElementLevel_Default,
			Title: titleForOther,
		},
		Divs:          cardCommon.Divs,
		ActionModules: cardCommon.ActionModules,
	}

	return cardForSelf, cardForOther, nil
}
