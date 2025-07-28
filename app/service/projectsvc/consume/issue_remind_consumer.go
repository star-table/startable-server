package consume

import (
	"fmt"
	"strings"

	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"github.com/star-table/startable-server/app/facade/msgfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errors"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/library/mq"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

func IssueRemindConsumer() {
	log.Infof("[IssueRemindConsumer] mq消息-任务提醒通知消费者启动成功")

	if config.GetMQ() == nil {
		log.Error("mq未配置")
		return
	}
	issueRemindConfig := config.GetMQ().Topics.IssueRemind

	if issueRemindConfig.Topic == "" {
		log.Error("[IssueRemindConsumer] mq issueRemind 未配置")
		return
	}

	client := *mq.GetMQClient()

	_ = client.ConsumeMessage(issueRemindConfig.Topic, issueRemindConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		//获取消息实体
		msgBo := &bo.IssueRemindMsg{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		log.Infof("[IssueRemindConsumer] mq消息-任务提醒通知-信息详情 %s", message.Body)
		// 任务提醒
		IssueRemind(msgBo)

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		msgBo := &bo.IssueRemindMqBo{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return
		}
		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, int(msgBo.PushType), 0)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})

}

func IssueRemind(msg *bo.IssueRemindMsg) {
	cardMsg := GetCardIssueBeOverDueRemind(msg)
	pushCard := &commonvo.PushCard{
		OrgId:         msg.OrgId,
		OutOrgId:      msg.OutOrgId,
		SourceChannel: msg.SourceChannel,
		OpenIds:       msg.OpenIds,
		CardMsg:       cardMsg,
	}

	errSys := card.PushCard(msg.OrgId, pushCard)
	if errSys != nil {
		log.Errorf("[IssueRemind] err:%v", errSys)
		return
	}
}

func GetCardIssueBeOverDueRemind(msg *bo.IssueRemindMsg) *pushV1.TemplateCard {
	title := msg.Title
	issueId := msg.IssueId
	orgId := msg.OrgId
	ownerId := msg.Operator
	sourceChannel := msg.SourceChannel

	issueLinks := domain.GetIssueLinks(sourceChannel, orgId, issueId)

	// 项目信息
	projectName := consts.CardDefaultIssueProjectName
	tableName := consts.DefaultTableName

	if msg.ProjectId > 0 {
		projectName = msg.ProjectName
	}
	if msg.TableId > 0 {
		tableName = msg.TableName
	}

	contentProject := ""
	if projectName != consts.CardDefaultIssueProjectName {
		contentProject = fmt.Sprintf(consts.CardTablePro, projectName, tableName)
	} else {
		contentProject = consts.CardDefaultIssueProjectName
	}

	parentTitle := consts.CardDefaultRelationIssueTitle
	if msg.ParentId > 0 {
		issueBos, err := domain.GetIssueInfosLc(orgId, 0, []int64{msg.ParentId})
		if err != nil {
			log.Errorf("[GetCardIssueBeOverDueRemind] domain.GetIssueInfosLc err:%v", err)
			return nil
		}
		if len(issueBos) <= 0 {
			return nil
		}
		parentTitle = issueBos[0].Title
	}

	// 负责人信息
	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: msg.OwnerId,
	})
	if userInfoBatch.Failure() {
		log.Errorf("[GetCardIssueBeOverDueRemind] orgfacade.GetBaseUserInfoBatch err:%v", userInfoBatch.Error())
		return nil
	}

	tableColumns, errSys := domain.GetTableColumnsMap(orgId, msg.TableId, []string{consts.BasicFieldTitle, consts.BasicFieldOwnerId})
	if errSys != nil {
		log.Errorf("[GetCardIssueBeOverDueRemind] GetTableColumnConfig err:%v, orgId:%v, issueId:%v", errSys, orgId, msg.IssueId)
		return nil
	}
	ownerSlice := []string{}
	for _, user := range userInfoBatch.BaseUserInfos {
		ownerSlice = append(ownerSlice, user.Name)
	}
	ownerDisplayName := strings.Join(ownerSlice, "，")
	if ownerDisplayName == "" {
		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	}

	cardMsg := &projectvo.CardBeOverdue{
		OrgId:          orgId,
		OwnerId:        ownerId,
		IssueId:        issueId,
		IssueLinks:     issueLinks,
		IssueTitle:     title,
		PlanEndTime:    msg.PlanEndTime,
		SourceChannel:  sourceChannel,
		AppId:          cast.ToString(msg.AppId),
		TableId:        cast.ToString(msg.TableId),
		ContentProject: contentProject,
		OwnerNameStr:   ownerDisplayName,
		ParentTitle:    parentTitle,
		HasPermission:  true,
		Tips:           "",
		TableColumn:    tableColumns,
	}
	cd := card.GetCardBeOverdue(cardMsg)
	log.Infof("[GetCardBeOverdue] card: %s", json.ToJsonIgnoreError(cd))
	return cd
}
