package notice

import (
	"strconv"
	"strings"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
)

func GetFeiShuAtConent(orgId int64, content string) string {
	userIds := util.GetCommentAtUserIds(content)
	//userId -> openId
	openIdMap := map[string]string{}
	if len(userIds) > 0 {
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
		if err != nil {
			log.Error(err)
		} else {
			for _, userInfo := range userInfos {
				openIdMap[strconv.FormatInt(userInfo.UserId, 10)] = userInfo.OutUserId
			}
		}
	}

	return util.RenderCommentContentToMarkDownWithOpenIdMap(content, false, openIdMap)
}

func GetFeiShuAtContentByUserIds(orgId int64, userIds []int64) (string, errs.SystemErrorInfo) {
	atSomeoneStr := strings.Builder{}
	if len(userIds) > 0 {
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
		if err != nil {
			log.Errorf("[GetFeiShuAtContentByUserIds] err: %v, orgId: %d", err, orgId)
			return atSomeoneStr.String(), err
		} else {
			for _, userInfo := range userInfos {
				atSomeoneStr.WriteString("<at id=" + userInfo.OutUserId + "></at>")
			}
		}
	}

	return atSomeoneStr.String(), nil
}
