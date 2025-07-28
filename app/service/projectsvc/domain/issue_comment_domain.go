package domain

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/app/facade/trendsfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
	"github.com/star-table/startable-server/common/model/vo/trendsvo"
)

func CreateIssueComment(issueBo bo.IssueBo, comment string, mentionedUserIds []int64, operatorId int64, attachmentIds []int64) (int64, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	projectId := issueBo.ProjectId

	//拼装评论
	commentBo := bo.CommentBo{
		OrgId:      orgId,
		ProjectId:  projectId,
		ObjectId:   issueBo.Id,
		ObjectType: consts.TrendsOperObjectTypeIssue,
		Content:    comment,
		Creator:    operatorId,
		Updator:    operatorId,
		IsDelete:   consts.AppIsNoDelete,
	}

	respVo := trendsfacade.CreateComment(trendsvo.CreateCommentReqVo{CommentBo: commentBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}

	commentId := respVo.CommentId
	commentBo.Id = commentId

	asyn.Execute(func() {
		pushType := consts.PushTypeIssueComment
		// 组装资源类型
		resourceInfos := []bo.ResourceInfoBo{}
		if attachmentIds != nil && len(attachmentIds) > 0 {
			isDelete := consts.AppIsNoDelete
			resourceBoList := resourcefacade.GetResourceBoList(resourcevo.GetResourceBoListReqVo{
				Page: 0,
				Size: 0,
				Input: resourcevo.GetResourceBoListCond{
					OrgId:       orgId,
					ResourceIds: &attachmentIds,
					IsDelete:    &isDelete,
				},
			})
			if resourceBoList.Failure() {
				log.Errorf("[CreateIssueComment] resourcefacade.GetResourceBoList failed, err:%v, orgId:%d, operatorId:%d, attachmentIds:%v",
					resourceBoList.Failure(), orgId, operatorId, attachmentIds)
				return
			}
			for _, resource := range *resourceBoList.ResourceBos {
				resourceInfos = append(resourceInfos, bo.ResourceInfoBo{
					Id:         resource.Id,
					Url:        resource.Host + resource.Path,
					Name:       resource.Name,
					Size:       resource.Size,
					UploadTime: resource.CreateTime,
					Suffix:     resource.Suffix,
				})
			}
		}

		issueTrendsBo := &bo.IssueTrendsBo{
			PushType:              pushType,
			OrgId:                 issueBo.OrgId,
			DataId:                issueBo.DataId,
			IssueId:               issueBo.Id,
			ParentIssueId:         issueBo.ParentId,
			ProjectId:             issueBo.ProjectId,
			TableId:               issueBo.TableId,
			PriorityId:            issueBo.PriorityId,
			ParentId:              issueBo.ParentId,
			OperatorId:            operatorId,
			IssueTitle:            issueBo.Title,
			IssueStatusId:         issueBo.Status,
			BeforeOwner:           issueBo.OwnerIdI64,
			AfterOwner:            issueBo.OwnerIdI64,
			BeforeChangeFollowers: issueBo.FollowerIdsI64,
			AfterChangeFollowers:  issueBo.FollowerIdsI64,
			Ext: bo.TrendExtensionBo{
				MentionedUserIds: mentionedUserIds,
				CommentBo:        &commentBo,
				ResourceInfo:     resourceInfos,
			},
		}

		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})

		asyn.Execute(func() {
			// 机器人推送个人卡片
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
		//推送群聊卡片
		orgRespVo := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
		if !orgRespVo.Failure() && orgRespVo.BaseOrgInfo.SourceChannel == sdk_const.SourceChannelFeishu {
			asyn.Execute(func() {
				PushInfoToChat(issueBo.OrgId, issueBo.ProjectId, issueTrendsBo, orgRespVo.BaseOrgInfo.SourceChannel)
			})
		}
	})

	return commentId, nil
}
