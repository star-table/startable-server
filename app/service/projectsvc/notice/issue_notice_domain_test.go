package notice

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang/request"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	"gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/orgfacade"

	fsSdkVo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/third_platform_sdk"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/smartystreets/goconvey/convey"
)

const (
	DingLinkWorkPlatform = "dingtalk://dingtalkclient/action/openapp?corpid=%s&container_type=work_platform&app_id=%s&redirect_type=jump&redirect_url=%s"
	DingLinkSlide        = "dingtalk://dingtalkclient/page/link?url=%s&pc_slide=true"
	DingAuth             = "%s/dd/auth?corpId=%s&redirect_url=%s"
)

func TestDingTalkCardLink(t *testing.T) {
	//https://app120908.eapps.dingtalkcloud.com/project/133884/task/1719260861299691520?projectTypeId=47&viewId=1719260861148696579&appId=1719260861148696578&issueId=5374112
	// dingtalk://dingtalkclient/action/openapp?corpid=企业的corpid&container_type=work_platform&app_id=appId&redirect_type=jump&redirect_url=跳转url
	convey.Convey("TestDingTalkCardLink", t, test.StartUp(func(ctx context.Context) {

		req := vo.GetIssueLinkReq{}
		req.Host = "https://app120908.eapps.dingtalkcloud.com"
		//req.Host = "https://app.startable.cn"
		req.CorpId = "ding2e2cd01df7535d08a1320dcb25e91351"
		req.ProjectId = 133884
		req.TableId = 1719260861299691520
		req.ProjectTypeId = 47
		req.AppId = 1719260861148696578
		req.IssueId = 5374112
		//redirectUri := fmt.Sprintf("/project/%d/task/%d?projectTypeId=%d&appId=%d&issueId=%d", req.ProjectId, req.TableId, req.ProjectTypeId, req.AppId,)
		redirectUri := fmt.Sprintf("/app/111/?projectTypeId=%d&issueId=%d&appId=%d", req.ProjectId, 0, req.ProjectTypeId, req.IssueId, req.AppId)
		noAuthLink := fmt.Sprintf("%s%s", req.Host, redirectUri)
		noAuthencodeUri := url.QueryEscape(noAuthLink)
		noauthUrlEncode := fmt.Sprintf(DingLinkWorkPlatform, req.CorpId, "120908", noAuthencodeUri)
		fmt.Sprintf("no auth ding url : %s", noauthUrlEncode)

		authLink := fmt.Sprintf(DingAuth, req.Host, req.CorpId, url.QueryEscape(redirectUri))
		encodeUri := url.QueryEscape(authLink)
		urlEncode := fmt.Sprintf(DingLinkWorkPlatform, req.CorpId, "120908", encodeUri)
		fmt.Sprintf("ding url : %s", urlEncode)
	}))
}

func GetTestTrendsBo() (bo.IssueTrendsBo, errs.SystemErrorInfo) {
	data := "{\"pushType\":1,\"orgId\":2373,\"operatorId\":29611,\"issueId\":10001356,\"parentIssueId\":0," +
		"\"projectId\":13908,\"priorityId\":0,\"parentId\":0,\"issueTitle\":\"su-p003-06141326-t001-2\",\"issueRemark\":\"\",\"issueStatusId\":16,\"issuePlanStartTime\":\"2022-06-27 17:22:29\",\"issuePlanEndTime\":\"2022-06-29 20:01:09\",\"sourceChannel\":\"fs\",\"beforeOwner\":[29611],\"afterOwner\":[29611],\"beforeChangeFollowers\":[29612],\"afterChangeFollowers\":[29612],\"beforeChangeParticipants\":null,\"afterChangeParticipants\":null,\"beforeChangeAuditors\":[],\"afterChangeAuditors\":null,\"beforeChangeRelating\":null,\"afterChangeRelating\":null,\"beforeChangeBaRelating\":null,\"afterChangeBaRelating\":null,\"beforeChangeDocuments\":null,\"afterChangeDocuments\":null,\"beforeChangeImages\":null,\"afterChangeImages\":null,\"beforeWorkHourIds\":null,\"afterWorkHourIds\":null,\"updateOwner\":false,\"updateFollower\":false,\"updateAuditor\":false,\"updateWorkHour\":false,\"issueChildren\":null,\"onlyNotice\":false,\"operateObjProperty\":\"\",\"newValue\":\"{\\\"createTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"ownerId\\\":[\\\"U_29611\\\"],\\\"ownerInfos\\\":null,\\\"projectTypeId\\\":0,\\\"propertyId\\\":0,\\\"planWorkHour\\\":0,\\\"status\\\":16,\\\"creator\\\":29611,\\\"dataId\\\":1536580934114050000,\\\"customField\\\":\\\"\\\",\\\"updator\\\":29611,\\\"typeForRelate\\\":0,\\\"issueStatusType\\\":2,\\\"id\\\":10001356,\\\"planStartTime\\\":\\\"2022-06-27 17:22:29\\\",\\\"iterationId\\\":0,\\\"moduleId\\\":0,\\\"remark\\\":\\\"\\\",\\\"issueProRelation\\\":null,\\\"path\\\":\\\"0,\\\",\\\"isFiling\\\":0,\\\"startTime\\\":\\\"2022-06-27 19:29:11\\\",\\\"auditStatus\\\":-1,\\\"versionId\\\":0,\\\"parentId\\\":0,\\\"orgId\\\":2373,\\\"code\\\":\\\"$13908-2\\\",\\\"title\\\":\\\"su-p003-06141326-t001-1\\\",\\\"planEndTime\\\":\\\"2022-06-29 20:01:09\\\",\\\"auditorIds\\\":[],\\\"followerInfos\\\":null,\\\"parentInfo\\\":null,\\\"resourceIds\\\":null,\\\"projectId\\\":13908,\\\"priorityId\\\":0,\\\"sourceId\\\":0,\\\"version\\\":1,\\\"auditorInfos\\\":null,\\\"lessData\\\":{\\\"priorityId\\\":0,\\\"delFlag\\\":2,\\\"version\\\":1,\\\"versionId\\\":0,\\\"orgId\\\":2373,\\\"startTime\\\":\\\"2022-06-27 19:29:11\\\",\\\"creator\\\":\\\"29611\\\",\\\"issueStatus\\\":16,\\\"appIds\\\":[\\\"1535463773907812353\\\"],\\\"remark\\\":\\\"\\\",\\\"projectId\\\":13908,\\\"followerIds\\\":[\\\"U_29612\\\"],\\\"parentId\\\":0,\\\"issueObjectTypeId\\\":0,\\\"ownerId\\\":[\\\"U_29611\\\"],\\\"createTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"planStartTime\\\":\\\"2022-06-27 17:22:29\\\",\\\"collaborators\\\":{\\\"auditorIds\\\":[],\\\"ownerId\\\":[\\\"U_29611\\\"],\\\"followerIds\\\":[\\\"U_29612\\\"]},\\\"title\\\":\\\"su-p003-06141326-t001\\\",\\\"path\\\":\\\"0,\\\",\\\"tableId\\\":\\\"1535463774599778304\\\",\\\"planWorkHour\\\":0,\\\"issueId\\\":10001356,\\\"projectObjectTypeId\\\":0,\\\"planEndTime\\\":\\\"2022-06-29 20:01:09\\\",\\\"ownerChangeTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"recycleFlag\\\":2,\\\"auditorIds\\\":[],\\\"sort\\\":\\\"655438865460\\\",\\\"id\\\":\\\"1536580934114050050\\\",\\\"auditStatus\\\":-1,\\\"order\\\":738721791,\\\"updateTime\\\":\\\"2022-06-28 17:22:31\\\",\\\"isDelete\\\":2,\\\"propertyId\\\":0,\\\"owner\\\":0,\\\"iterationId\\\":0,\\\"status\\\":1,\\\"workHour\\\":{\\\"actualHour\\\":\\\"0\\\",\\\"planHour\\\":\\\"1.2\\\",\\\"collaboratorIds\\\":[\\\"U_29611\\\"]},\\\"updator\\\":\\\"29611\\\",\\\"sourceId\\\":0,\\\"isFiling\\\":0,\\\"moduleId\\\":0,\\\"code\\\":\\\"$13908-2\\\",\\\"issueStatusType\\\":2,\\\"endTime\\\":\\\"1970-01-01 00:00:00\\\"},\\\"followerIds\\\":[\\\"U_29612\\\"],\\\"participantInfos\\\":null,\\\"sort\\\":655438865460,\\\"updateTime\\\":\\\"2022-06-28 17:22:31\\\",\\\"isDelete\\\":2,\\\"tags\\\":null,\\\"projectObjectTypeId\\\":0,\\\"ownerChangeTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"issueObjectTypeId\\\":0,\\\"endTime\\\":\\\"1970-01-01 00:00:00\\\",\\\"parentTitle\\\":\\\"\\\",\\\"tableId\\\":1535463774599778300,\\\"appId\\\":1535463773907812400}\",\"oldValue\":\"{\\\"participantInfos\\\":null,\\\"parentTitle\\\":\\\"\\\",\\\"recycleFlag\\\":2,\\\"id\\\":10001356,\\\"version\\\":1,\\\"customField\\\":\\\"\\\",\\\"moduleId\\\":0,\\\"ownerChangeTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"priorityId\\\":0,\\\"appId\\\":1535463773907812400,\\\"createTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"collaborators\\\":{\\\"followerIds\\\":[\\\"U_29612\\\"],\\\"auditorIds\\\":[],\\\"ownerId\\\":[\\\"U_29611\\\"]},\\\"parentInfo\\\":null,\\\"auditorInfos\\\":null,\\\"issueStatusType\\\":2,\\\"orgId\\\":2373,\\\"ownerId\\\":[\\\"U_29611\\\"],\\\"ownerInfos\\\":null,\\\"lessData\\\":{\\\"isDelete\\\":2,\\\"path\\\":\\\"0,\\\",\\\"isFiling\\\":0,\\\"startTime\\\":\\\"2022-06-27 19:29:11\\\",\\\"planEndTime\\\":\\\"2022-06-29 20:01:09\\\",\\\"endTime\\\":\\\"1970-01-01 00:00:00\\\",\\\"delFlag\\\":2,\\\"recycleFlag\\\":2,\\\"createTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"collaborators\\\":{\\\"followerIds\\\":[\\\"U_29612\\\"],\\\"auditorIds\\\":[],\\\"ownerId\\\":[\\\"U_29611\\\"]},\\\"propertyId\\\":0,\\\"ownerChangeTime\\\":\\\"2022-06-14 13:26:57\\\",\\\"sort\\\":\\\"655438865460\\\",\\\"order\\\":738721791,\\\"owner\\\":0,\\\"auditorIds\\\":[],\\\"issueStatus\\\":16,\\\"auditStatus\\\":-1,\\\"status\\\":1,\\\"title\\\":\\\"su-p003-06141326-t001\\\",\\\"moduleId\\\":0,\\\"issueObjectTypeId\\\":0,\\\"ownerId\\\":[\\\"U_29611\\\"],\\\"parentId\\\":0,\\\"creator\\\":\\\"29611\\\",\\\"version\\\":1,\\\"updator\\\":\\\"29611\\\",\\\"issueStatusType\\\":2,\\\"followerIds\\\":[\\\"U_29612\\\"],\\\"projectId\\\":13908,\\\"id\\\":\\\"1536580934114050050\\\",\\\"sourceId\\\":0,\\\"orgId\\\":2373,\\\"workHour\\\":{\\\"collaboratorIds\\\":[\\\"U_29611\\\"],\\\"actualHour\\\":\\\"0\\\",\\\"planHour\\\":\\\"1.2\\\"},\\\"iterationId\\\":0,\\\"tableId\\\":\\\"1535463774599778304\\\",\\\"code\\\":\\\"$13908-2\\\",\\\"issueId\\\":10001356,\\\"updateTime\\\":\\\"2022-06-28 17:22:31\\\",\\\"versionId\\\":0,\\\"planWorkHour\\\":0,\\\"planStartTime\\\":\\\"2022-06-27 17:22:29\\\",\\\"remark\\\":\\\"\\\",\\\"priorityId\\\":0,\\\"appIds\\\":[\\\"1535463773907812353\\\"],\\\"projectObjectTypeId\\\":0},\\\"issueId\\\":10001356,\\\"planEndTime\\\":\\\"2022-06-29 20:01:09\\\",\\\"startTime\\\":\\\"2022-06-27 19:29:11\\\",\\\"workHour\\\":{\\\"collaboratorIds\\\":[\\\"U_29611\\\"],\\\"actualHour\\\":\\\"0\\\",\\\"planHour\\\":\\\"1.2\\\"},\\\"projectId\\\":13908,\\\"typeForRelate\\\":0,\\\"resourceIds\\\":null,\\\"code\\\":\\\"$13908-2\\\",\\\"planStartTime\\\":\\\"2022-06-27 17:22:29\\\",\\\"updator\\\":29611,\\\"updateTime\\\":\\\"2022-06-28 17:22:31\\\",\\\"dataId\\\":1536580934114050000,\\\"order\\\":738721791,\\\"appIds\\\":[\\\"1535463773907812353\\\"],\\\"delFlag\\\":2,\\\"path\\\":\\\"0,\\\",\\\"isFiling\\\":0,\\\"issueObjectTypeId\\\":0,\\\"planWorkHour\\\":0,\\\"isDelete\\\":2,\\\"followerInfos\\\":null,\\\"projectTypeId\\\":0,\\\"projectObjectTypeId\\\":0,\\\"title\\\":\\\"su-p003-06141326-t001\\\",\\\"issueProRelation\\\":null,\\\"parentId\\\":0,\\\"followerIds\\\":[\\\"U_29612\\\"],\\\"auditStatus\\\":-1,\\\"tags\\\":null,\\\"propertyId\\\":0,\\\"versionId\\\":0,\\\"status\\\":16,\\\"remark\\\":\\\"\\\",\\\"owner\\\":0,\\\"sourceId\\\":0,\\\"endTime\\\":\\\"1970-01-01 00:00:00\\\",\\\"issueStatus\\\":16,\\\"sort\\\":655438865460,\\\"auditorIds\\\":[],\\\"tableId\\\":1535463774599778300,\\\"iterationId\\\":0,\\\"creator\\\":29611}\",\"ext\":{\"issueType\":\"T\",\"objName\":\"su-p003-06141326-t001\",\"changeList\":[{\"field\":\"title\",\"fieldType\":\"input\",\"fieldName\":\"标题\",\"fieldNameValue\":\"\",\"aliasTitle\":\"\",\"oldValue\":\"su-p003-06141326-t001\",\"newValue\":\"su-p003-06141326-t001-1\",\"oldUserIdsOrDeptIdsValue\":null,\"newUserIdsOrDeptIdsValue\":null,\"changeTypeDesc\":\"\",\"isForWorkHour\":false}],\"memberInfo\":null,\"tagInfo\":null,\"relationIssue\":{\"id\":0,\"title\":\"\"},\"commonChange\":null,\"folderId\":0,\"mentionedUserIds\":null,\"commentBo\":{\"id\":0,\"orgId\":0,\"projectId\":0,\"trendsId\":0,\"objectId\":0,\"objectType\":\"\",\"content\":\"\",\"parentId\":0,\"creator\":0,\"createTime\":\"0001-01-01 00:00:00\",\"updator\":0,\"updateTime\":\"0001-01-01 00:00:00\",\"version\":0,\"isDelete\":0},\"resourceInfo\":null,\"remark\":\"\",\"fieldIds\":null,\"projectObjectTypeId\":0,\"auditInfo\":{\"status\":0,\"remark\":\"\",\"attachments\":null},\"formConfig\":{\"fields\":null,\"fieldOrders\":null,\"viewOrders\":null,\"baseFields\":null},\"addedFormFields\":null,\"deletedFormFields\":null,\"updatedFormFields\":null},\"operateTime\":\"0001-01-01 00:00:00\",\"tableId\":1535463774599778304}"
	trendsBo := bo.IssueTrendsBo{}
	if err := json.FromJson(data, &trendsBo); err != nil {
		log.Errorf("[GetTestTrendsBo] err: %v", err)
		return trendsBo, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	return trendsBo, nil
}

func TestCardPush1(t *testing.T) {
	convey.Convey("TestCardPush", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)

		title := "KZ专发测试专用"
		meta := &commonvo.CardMeta{}
		meta.IsWide = false
		meta.Title = title
		div := &commonvo.CardDiv{}
		div.Fields = append(div.Fields, &commonvo.CardField{
			Key:   "已逾期",
			Value: "381",
		})
		div.Fields = append(div.Fields, &commonvo.CardField{
			Key:   "今日截止",
			Value: "01",
		})
		div.Fields = append(div.Fields, &commonvo.CardField{
			Key:   "即将逾期",
			Value: "0",
		})
		div.Fields = append(div.Fields, &commonvo.CardField{
			Key:   "待完成",
			Value: "129",
		})
		div.Fields = append(div.Fields, &commonvo.CardField{
			Key:   "项目/表",
			Value: "129",
		})
		meta.Divs = append(meta.Divs, div)

		url := "http://baidu.com"
		meta.ActionMarkdowns = append(meta.ActionMarkdowns, fmt.Sprintf(consts.MarkdownLink, consts.CardButtonTextForViewDetail, url))
		meta.FsActionElements = []interface{}{
			fsSdkVo.CardElementActionModule{
				Tag:    "action",
				Layout: "bisected",
				Actions: []interface{}{
					fsSdkVo.ActionButton{
						Tag: "button",
						Text: fsSdkVo.CardElementText{
							Tag:     "plain_text",
							Content: consts.CardButtonTextForViewDetail,
						},
						Url:  url,
						Type: consts.FsCardButtonColorPrimary,
					},
					fsSdkVo.ActionButton{
						Tag: "button",
						Text: fsSdkVo.CardElementText{
							Tag:     "plain_text",
							Content: consts.CardButtonTextForViewInsideApp,
						},
						Url:  url,
						Type: consts.FsCardButtonColorDefault,
					},
				},
			},
			fsSdkVo.CardElementActionModule{
				Tag: "hr",
			},
			fsSdkVo.CardElementActionModule{
				Tag: "action",
				Actions: []interface{}{
					fsSdkVo.ActionButton{
						Tag: "button",
						Text: fsSdkVo.CardElementText{
							Tag:     "plain_text",
							Content: "点击退订",
						},
						Type: "link",
						Value: map[string]interface{}{
							consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
							consts.FsCardValueAction:   consts.FsCardUnsubscribeIssueRemind,
							consts.FsCardValueOrgId:    0,
						},
					},
				},
			},
		}

		card.SendCardMeta(sdk_const.SourceChannelFeishu, "12a813fb180f1758", meta, []string{"ou_dcf37513afc11d505b80e080c84c5aa2", "ou_5d296a1582a335dd5d9ddafa694a9e2d", "ou_1e82ae07440dfcee9f9dad91483b0be9"})
		card.SendCardMeta(sdk_const.SourceChannelWeixin, "wwf36b5e6ef0b569ac", meta, []string{"bac67275e863edd81145365b953dd581", "WuKuai", "puck"})
		card.SendCardMeta(sdk_const.SourceChannelDingTalk, "dingae89bb4c8cbaa471a39a90f97fcb1e09", meta, []string{"36203862793986", "01444707491231132620", "01082729076720711766"})
	}))
}

func TestPushIssueByChannel(t *testing.T) {
	trendsBo, _ := GetTestTrendsBo()
	convey.Convey("TestPushIssueByChannel", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		meta := &commonvo.CardMeta{}
		str := `{
    "IsWide": false,
    "Level": 0,
    "Title": "苏汉宇 评论了记录",
    "Divs": [
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "苏汉宇：",
                    "Value": "\t\t\t“<at id=ou_2bbc27672c6ce4f0c2c73856f73ccbdb></at>说什么呢”"
                }
            ]
        }
    ],
    "ActionMarkdowns": [
        "[查看详情](https://applink.feishu.cn/client/web_app/open?appId=cli_a02f641485b9d00b&path=/feishu/auth?callback=L3Byb2plY3QvNjA5NDQvdGFzay8xNTYxOTc5NjY1NzE3OTIzODQwP3Byb2plY3RUeXBlSWQ9NDcmYXBwSWQ9MTU2MDUxMzA5NzA2MjUxNDY4OSZpc3N1ZUlkPTEwMDg0NzM3)"
    ],
    "FsActionElements": [
        {
            "tag": "action",
            "layout": "bisected",
            "actions": [
                {
                    "tag": "button",
                    "text": {
                        "tag": "plain_text",
                        "content": "查看详情"
                    },
                    "url": "https://applink.feishu.cn/client/web_app/open?appId=cli_a02f641485b9d00b&mode=sidebar-semi&path=/feishu/auth&callback=L3Rhc2tEZXRhaWw_aXNzdWVJZD0xMDA4NDczNw$$",
                    "type": "primary"
                },
                {
                    "tag": "button",
                    "text": {
                        "tag": "plain_text",
                        "content": "应用内查看"
                    },
                    "url": "https://applink.feishu.cn/client/web_app/open?appId=cli_a02f641485b9d00b&path=/feishu/auth?callback=L3Byb2plY3QvNjA5NDQvdGFzay8xNTYxOTc5NjY1NzE3OTIzODQwP3Byb2plY3RUeXBlSWQ9NDcmYXBwSWQ9MTU2MDUxMzA5NzA2MjUxNDY4OSZpc3N1ZUlkPTEwMDg0NzM3",
                    "type": "default"
                }
            ]
        },
        {
            "tag": "hr",
            "actions": null
        },
        {
            "tag": "action",
            "actions": [
                {
                    "tag": "button",
                    "text": {
                        "tag": "plain_text",
                        "content": "点击退订"
                    },
                    "type": "link",
                    "value": {
                        "cardType": "NoticeSubscribe",
                        "action": "FsCardActionUnsubscribeIssueRemind",
                        "orgId": 2373
                    }
                }
            ]
        }
    ]
}`
		json.FromJson(str, meta)
		card.GenerateFeiShuCard(meta)
		PushIssueByChannel(trendsBo, sdk_const.SourceChannelFeishu)
	}))
}

func TestIssueChatCardPush1(t *testing.T) {
	trendsBo, _ := GetTestTrendsBo()
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		if err := domain.PushInfoToChat(2373, 13908, &trendsBo, "fs"); err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}

//func TestProjectMemberChangeNotice2(t *testing.T) {
//	convey.Convey("TestPushIssueByChannel", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
//		ProjectMemberChangeNotice(bo.ProjectMemberChangeBo{
//			PushType:            1002,
//			OrgId:               1113,
//			ProjectId:           1379,
//			OperatorId:          1293,
//			BeforeChangeMembers: nil,
//			AfterChangeMembers:  []int64{1289},
//			OperateObjProperty:  "",
//			NewValue:            "",
//			OldValue:            "",
//		})
//	}))
//}

func TestGetNormalUserIdsWithFilter(t *testing.T) {
	convey.Convey("TestPushIssueByChannel", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		fmt.Println(GetNormalUserIdsWithFilter(1242, 1882, []int64{1882}, []int64{}, []int64{1883}, []int64{}, "fs",
			1, nil))
	}))
}

func TestPushIssueComment(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		str := `{
    "pushType": 9,
    "orgId": 2373,
    "operatorId": 29611,
    "issueId": 10085511,
    "parentIssueId": 0,
    "projectId": 60944,
    "priorityId": 0,
    "parentId": 0,
    "issueTitle": "哈哈1111-move111",
    "issueRemark": "",
    "issueStatusId": 7,
    "issuePlanStartTime": null,
    "issuePlanEndTime": null,
    "sourceChannel": "",
    "beforeOwner": [
        29611
    ],
    "afterOwner": [
        29611
    ],
    "beforeChangeFollowers": [
        29611
    ],
    "afterChangeFollowers": [
        29611
    ],
    "beforeChangeParticipants": null,
    "afterChangeParticipants": null,
    "beforeChangeAuditors": null,
    "afterChangeAuditors": null,
    "beforeChangeRelating": null,
    "afterChangeRelating": null,
    "beforeChangeBaRelating": null,
    "afterChangeBaRelating": null,
    "beforeChangeDocuments": null,
    "afterChangeDocuments": null,
    "beforeChangeImages": null,
    "afterChangeImages": null,
    "beforeWorkHourIds": null,
    "afterWorkHourIds": null,
    "updateOwner": false,
    "updateFollower": false,
    "updateAuditor": false,
    "updateWorkHour": false,
    "issueChildren": null,
    "onlyNotice": false,
    "operateObjProperty": "",
    "newValue": "",
    "oldValue": "",
    "ext": {
        "issueType": "",
        "objName": "",
        "changeList": null,
        "memberInfo": null,
        "tagInfo": null,
        "relationIssue": {
            "id": 0,
            "title": ""
        },
        "commonChange": null,
        "folderId": 0,
        "mentionedUserIds": [
            29612,
            33988
        ],
        "commentBo": {
            "id": 3572,
            "orgId": 2373,
            "projectId": 60944,
            "trendsId": 0,
            "objectId": 10085511,
            "objectType": "Issue",
            "content": "@#[成卫忠1:29612]&$@#[冯建:33988]&$就好撒",
            "parentId": 0,
            "creator": 29611,
            "createTime": "0001-01-01 00:00:00",
            "updator": 29611,
            "updateTime": "0001-01-01 00:00:00",
            "version": 0,
            "isDelete": 2
        },
        "resourceInfo": [],
        "remark": "",
        "fieldIds": null,
        "projectObjectTypeId": 0,
        "auditInfo": {
            "status": 0,
            "remark": "",
            "attachments": null
        },
        "formConfig": {
            "fields": null,
            "fieldOrders": null,
            "viewOrders": null,
            "baseFields": null
        },
        "addedFormFields": null,
        "deletedFormFields": null,
        "updatedFormFields": null
    },
    "operateTime": "2022-08-31 19:58:56",
    "tableId": 1561997815561850880
}`
		issueTrendsBo := &bo.IssueTrendsBo{}
		issueNoticeBo := &bo.IssueNoticeBo{}
		json.FromJson(str, issueTrendsBo)
		copyer.Copy(issueTrendsBo, issueNoticeBo)
		ext := issueTrendsBo.Ext
		content := ext.CommentBo.Content
		mentionedUserIds := ext.MentionedUserIds
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		PushIssueComment(*issueTrendsBo, content, mentionedUserIds, issueNoticeBo.PushType)
	}))
}

func TestPushIssue(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		str := `{
    "pushType": 0,
    "orgId": 2373,
    "operatorId": 29610,
    "issueId": 10100354,
    "parentIssueId": 0,
    "projectId": 14207,
    "priorityId": 0,
    "parentId": 0,
    "issueTitle": "77777",
    "issueRemark": "",
    "issueStatusId": 7,
    "issuePlanStartTime": "0001-01-01 00:00:00",
    "issuePlanEndTime": "0001-01-01 00:00:00",
    "sourceChannel": "",
    "beforeOwner": [
        29610
    ],
    "afterOwner": [
        29610
    ],
    "beforeChangeFollowers": [
        29610
    ],
    "afterChangeFollowers": [
        29610
    ],
    "beforeChangeParticipants": null,
    "afterChangeParticipants": null,
    "beforeChangeAuditors": [],
    "afterChangeAuditors": [],
    "beforeChangeRelating": {
        "linkTo": null,
        "linkFrom": null
    },
    "afterChangeRelating": {
        "linkTo": null,
        "linkFrom": null
    },
    "beforeChangeBaRelating": {
        "linkTo": null,
        "linkFrom": null
    },
    "afterChangeBaRelating": {
        "linkTo": null,
        "linkFrom": null
    },
    "beforeChangeDocuments": null,
    "afterChangeDocuments": null,
    "beforeChangeImages": null,
    "afterChangeImages": null,
    "beforeWorkHourIds": null,
    "afterWorkHourIds": null,
    "updateOwner": false,
    "updateFollower": false,
    "updateAuditor": false,
    "updateWorkHour": false,
    "issueChildren": null,
    "onlyNotice": false,
    "operateObjProperty": "",
    "newValue": "{\"issueStatus\":7,\"iterationId\":0,\"code\":\"$14207-38\",\"parentId\":0,\"tableId\":\"1549934550950350848\",\"followerIds\":[\"U_29610\"],\"title\":\"77777\",\"ownerId\":[\"U_29610\"],\"path\":\"0,\",\"appIds\":[\"1549934550279315457\"],\"recycleFlag\":2,\"delFlag\":2,\"ownerChangeTime\":\"2022-10-10 20:09:13\",\"projectId\":14207,\"auditStatusDetail\":{},\"creator\":\"29610\",\"createTime\":\"2022-10-10 20:09:13\",\"orgId\":2373,\"order\":800194559,\"updator\":\"29610\",\"updateTime\":\"2022-10-10 20:09:13\",\"issueStatusType\":1,\"auditStatus\":-1,\"issueId\":10100354}",
    "oldValue": "",
    "ext": {
        "issueType": "",
        "objName": "",
        "changeList": null,
        "memberInfo": null,
        "tagInfo": null,
        "relationIssue": {
            "id": 0,
            "title": ""
        },
        "commonChange": null,
        "folderId": 0,
        "mentionedUserIds": null,
        "commentBo": {
            "id": 0,
            "orgId": 0,
            "projectId": 0,
            "trendsId": 0,
            "objectId": 0,
            "objectType": "",
            "content": "",
            "parentId": 0,
            "creator": 0,
            "createTime": "0001-01-01 00:00:00",
            "updator": 0,
            "updateTime": "0001-01-01 00:00:00",
            "version": 0,
            "isDelete": 0
        },
        "resourceInfo": null,
        "remark": "",
        "fieldIds": null,
        "projectObjectTypeId": 0,
        "auditInfo": {
            "status": 0,
            "remark": "",
            "attachments": null
        },
        "formConfig": {
            "fields": null,
            "fieldOrders": null,
            "viewOrders": null,
            "baseFields": null
        },
        "addedFormFields": null,
        "deletedFormFields": null,
        "updatedFormFields": null
    },
    "operateTime": "2022-10-10 20:09:13",
    "tableId": 1549934550950350848
}`
		issueTrendsBo := &bo.IssueTrendsBo{}
		issueNoticeBo := &bo.IssueNoticeBo{}
		json.FromJson(str, issueTrendsBo)
		copyer.Copy(issueTrendsBo, issueNoticeBo)
		//ext := issueTrendsBo.Ext
		//content := ext.CommentBo.Content
		//mentionedUserIds := ext.MentionedUserIds
		//PushIssueComment(*issueTrendsBo, content, mentionedUserIds, issueNoticeBo.PushType)
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		PushIssue(*issueTrendsBo)
	}))
}

func TestSendCardMetaDing(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		str := `{
    "IsWide": false,
    "Level": 0,
    "Title": "苏汉宇 更新了记录",
    "Divs": [
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "标题",
                    "Value": "saashb地方看看"
                }
            ]
        },
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "项目/表",
                    "Value": "111/任务"
                }
            ]
        },
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "负责人",
                    "Value": "成卫忠"
                }
            ]
        },
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "任务状态",
                    "Value": "~~未开始~~ **→** 进行中"
                }
            ]
        }
    ],
    "ActionMarkdowns": [
        "[查看详情](dingtalk://dingtalkclient/action/openapp?corpid=noNeedTokenCorpId&container_type=work_platform&app_id=75864&redirect_type=jump&redirect_url=%2Fproject%2F60960%2Ftask%2F1560571170271531008%3FprojectTypeId%3D47%26appId%3D1560571169839611905%26issueId%3D10084722)"
    ],
    "FsActionElements": [
        {
            "tag": "action",
            "layout": "bisected",
            "actions": [
                {
                    "tag": "button",
                    "text": {
                        "tag": "plain_text",
                        "content": "查看详情"
                    },
                    "url": "http://www.baidu.com",
                    "type": "primary"
                },
                {
                    "tag": "button",
                    "text": {
                        "tag": "plain_text",
                        "content": "应用内查看"
                    },
                    "url": "dingtalk://dingtalkclient/action/openapp?corpid=noNeedTokenCorpId&container_type=work_platform&app_id=75864&redirect_type=jump&redirect_url=%2Fproject%2F60960%2Ftask%2F1560571170271531008%3FprojectTypeId%3D47%26appId%3D1560571169839611905%26issueId%3D10084722",
                    "type": "default"
                }
            ]
        },
        {
            "tag": "hr",
            "actions": null
        },
        {
            "tag": "action",
            "actions": [
                {
                    "tag": "button",
                    "text": {
                        "tag": "plain_text",
                        "content": "点击退订"
                    },
                    "type": "link",
                    "value": {
                        "cardType": "NoticeSubscribe",
                        "action": "FsCardActionUnsubscribeIssueRemind",
                        "orgId": 2588
                    }
                }
            ]
        }
    ]
}`
		corpId := "dingae89bb4c8cbaa471a39a90f97fcb1e09"
		meta := &commonvo.CardMeta{}
		json.FromJson(str, meta)
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		//card.SendCardMetaFeiShu(corpId, meta, []string{"ou_3ab7fe596cf91692218f744558ae157f"})
		card.SendCardMetaDingTalk(corpId, meta, []string{"01023369511624811493"})
	}))
}

func TestSendCardMetaFs(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		str := `{
    "IsWide": false,
    "Level": 1,
    "Title": "⏰ 即将到达截止时间",
    "Divs": [
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "标题",
                    "Value": "与之都行"
                }
            ]
        },
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "项目/表",
                    "Value": "cwz-催办测试/任务"
                }
            ]
        },
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "负责人",
                    "Value": "苏汉宇，成卫忠1"
                }
            ]
        },
        {
            "Fields": [
                {
                    "Level": 0,
                    "Key": "计划截止时间",
                    "Value": ""
                }
            ]
        }
    ],
    "ActionMarkdowns": [
        "[查看详情](https://applink.feishu.cn/client/web_app/open?appId=cli_a02f641485b9d00b&path=/feishu/auth?callback=L3Byb2plY3QvMTM5NDUvdGFzay8xNTM3Njc0OTk1MTA1MjcxODA4P3Byb2plY3RUeXBlSWQ9MSZhcHBJZD0xNTM3Njc0OTk0NDE3NTAwMTYxJmlzc3VlSWQ9MTAwMTM4NTA$)"
    ]
}`
		meta := &commonvo.CardMeta{}
		json.FromJson(str, meta)
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		card.SendCardMetaFeiShu("12a813fb180f1758", meta, []string{"ou_2bbc27672c6ce4f0c2c73856f73ccbdb"})
	}))
}

func TestGetIssueLinks(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		domain.GetIssueLinks("weixin", 2673, 10085065)
	}))
}

func TestGetSKUPage(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		sdk, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, "dingae89bb4c8cbaa471a39a90f97fcb1e09")
		if err != nil {
			t.Error(err)
		}
		ding := sdk.GetOriginClient()
		dingOriginClient := ding.(*dingtalk.DingTalk)
		resp, err := dingOriginClient.GetOrderSkuPage(&request.GetSkuPageReq{
			GoodsCode: "DT_GOODS_881663135459256",
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(resp.Result)
	}))
}
