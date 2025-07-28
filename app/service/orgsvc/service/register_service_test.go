package orgsvc

import (
	"testing"

	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/smartystreets/goconvey/convey"
)

func TestUserRegister(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		authCode := "000000"
		name := "abc"
		t.Log(UserRegister(orgvo.UserRegisterReqVo{Input: vo.UserRegisterReq{
			SourceChannel:  "fs",
			SourcePlatform: "fs",
			AuthCode:       &authCode,
			Name:           &name,
			UserName:       "17700000002",
			RegisterType:   1,
		}}))
	}))
}
