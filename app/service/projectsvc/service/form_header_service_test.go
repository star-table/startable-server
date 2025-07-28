package service

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/core/util/json"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/model/vo/orgvo"

	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/projectvo"

	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetFormConfigBatch(t *testing.T) {
	//convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
	//	res, err := GetFormConfigBatch(projectvo.GetFormConfigBatchReq{
	//		OrgId:  1629,
	//		UserId: 0,
	//		AppIds: []int64{1452546873305956354},
	//	})
	//	t.Log(json.ToJsonIgnoreError(res))
	//	t.Log(err)
	//}))
}

func TestUpdateColumnTriggerCollaborate1(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		err := domain.DealCollaboratorsColumnsToIssueChat(2373, 1535463774599778304, &projectvo.TableColumnData{
			Name:              consts.BasicFieldFollowerIds,
			Label:             "关注人",
			Writable:          true,
			Editable:          true,
			Unique:            false,
			UniquePreHandler:  "",
			SensitiveStrategy: "",
			SensitiveFlag:     0,
			Field: projectvo.TableColumnField{
				Type:       consts.LcColumnFieldTypeMember,
				CustomType: "",
				DataType:   "",
				Props: map[string]interface{}{
					"multiple":          false,
					"limit":             nil,
					"collaboratorRoles": nil,
					"member": lc_table.LcPropMember{
						Multiple:        false,
						Required:        false,
						IsCollaborators: true,
					},
				},
			},
		},
			consts.UpdateColumn)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}

func TestGetUsersByOpenIdBatch1(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		joinUsersResp := orgfacade.GetBaseUserInfoByEmpIdBatch(orgvo.GetBaseUserInfoByEmpIdBatchReqVo{
			OrgId: 2373,
			Input: orgvo.GetBaseUserInfoByEmpIdBatchReqVoInput{
				OpenIds: []string{"ou_3ab7fe596cf91692218f744558ae157f"},
			},
		})
		if joinUsersResp.Failure() {
			log.Errorf("[ChatUserJoinHandler.Handle] GetBaseUserInfoByEmpIdBatch err: %v", joinUsersResp.Error())
			return
		}
		t.Log(json.ToJsonIgnoreError(joinUsersResp))
	}))

}
