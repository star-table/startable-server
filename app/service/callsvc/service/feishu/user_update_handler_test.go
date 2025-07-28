package callsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

//func TestUserUpdateHandler_Handle(t *testing.T) {
//	convey.Convey("Test TestPushBotHelpNotice", t, test.StartUp(func(ctx context.Context) {
//		handler := &UserUpdateHandler{}
//	}))
//}

func TestUserUpdateHandler_Handle(t *testing.T) {
	convey.Convey("Test TestPushBotHelpNotice", t, test.StartUp(func(ctx context.Context) {
		//orgId := int64(1004)
		//// cache 有值时，同步部分企业的员工信息
		//// cache 无值时，同步所有企业的员工信息 polaris:orgsvc:sys:fs_user_update_white_list
		//key := sconsts.CacheFsUserUpdateWhiteList
		//value, err := cache.Get(key)
		//if err != nil {
		//	log.Error(err)
		//	return
		//}
		//if value != "" {
		//	fsUserOrgIdWhiteList := &[]int64{}
		//	err = json.FromJson(value, fsUserOrgIdWhiteList)
		//	if err != nil {
		//		log.Info(strs.ObjectToString(err))
		//		return
		//	}
		//	// 当前企业不再白名单当中，则不处理
		//	if ok, _ := slice.Contain(*fsUserOrgIdWhiteList, orgId); !ok {
		//		return
		//	}
		//}
		//result := &[]int64{}
		//err = json.FromJson(value, result)
		//if err != nil {
		//	log.Info(strs.ObjectToString(err))
		//	return
		//}
		//t.Logf("end...\n")
		jsonStr := "{\"name\":\"a\",\"age\":1}"
		type A struct {
			Name string `json:"name"`
			Age  int64  `json:"age"`
		}
		res := &A{}
		err := json.FromJson(jsonStr, res)
		t.Log(*res)
		t.Log(err)

	}))
}
