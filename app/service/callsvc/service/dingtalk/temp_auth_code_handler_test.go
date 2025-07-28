package callsvc

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/test"
	_ "github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
)

//func TestGetDepUser(t *testing.T) {
//
//	convey.Convey("Test config", t, test.StartUp(func(ctx context.Context) {
//
//		client, _ := dingtalk.GetDingTalkClientRest("ding8ac2bab2b708b3cc35c2f4657eb6378f")
//
//		userListResp, _ := client.GetDepMemberDetailList("0", "zh_CN", 1, 100, "")
//
//		t.Log(userListResp.UserList)
//		fetchChild := false
//		resp, _ := client.GetDeptList(nil, &fetchChild, "1")
//		t.Log(json.ToJson(resp))
//		t.Log(len(resp.Department))
//
//		if len(resp.Department) > 0 {
//
//			userIds := &[]string{}
//			for _, dept := range resp.Department {
//				userListResp, _ := client.GetDepMemberIds(strconv.FormatInt(dept.Id, 10))
//
//				if len(userListResp.UserIds) > 0 {
//					*userIds = append(*userIds, userListResp.UserIds...)
//				}
//			}
//			userListResp, _ := client.GetDepMemberIds("1")
//
//			if len(userListResp.UserIds) > 0 {
//				*userIds = append(*userIds, userListResp.UserIds...)
//			}
//
//			*userIds = slice.SliceUniqueString(*userIds)
//
//			t.Log(len(*userIds))
//			t.Log(json.ToJson(userIds))
//		}
//
//	}))
//}
//
//func TestDingTalkGetSuiteTicket(t *testing.T) {
//
//	convey.Convey("Test TestDingTalkGetSuiteTicket", t, test.StartUp(func(ctx context.Context) {
//
//		suiteTicket, err := dingtalk.GetSuiteTicket()
//
//		if err != nil {
//			fmt.Println("获取suiteTicket异常")
//			return
//		}
//
//		fmt.Println("钉钉推送的suiteTicket  =", suiteTicket)
//
//	}))
//}

//func TestInitDepartment(t *testing.T) {
//
//	convey.Convey("Test InitDepartment", t, func() {
//
//		//开启事务
//		conn, err := mysql.GetConnect()
//		if err != nil {
//			log.Errorf(consts.DBOpenErrorSentence, err)
//		}
//		tx, err := conn.NewTx(nil)
//		if err != nil {
//			log.Errorf(consts.TxOpenErrorSentence, err)
//		}
//		defer mysql.Close(conn, tx)
//
//		rand.Seed(date.Now().Unix())
//		intn := rand.Intn(10000)
//
//		err = initDepartment(int64(intn), tx, "ding8ac2bab2b708b3cc35c2f4657eb6378f")
//		if err != nil {
//			tx.Rollback()
//			fmt.Println(err)
//		}
//		tx.Commit()
//	})
//}

//func TestInitUserRoles(t *testing.T) {
//
//	convey.Convey("Test InitUserRoles", t, test.StartUp(func(ctx context.Context) {
//
//		//开启事务
//		conn, err := mysql.GetConnect()
//		if err != nil {
//			log.Errorf(consts.DBOpenErrorSentence, err)
//		}
//		tx, err := conn.NewTx(nil)
//		if err != nil {
//			log.Errorf(consts.TxOpenErrorSentence, err)
//		}
//		defer mysql.Close(conn, tx)
//
//		convey.Convey("mock InitUserRoles", func() {
//
//			//TODO  这边需要自己填写一个corpId不然 projectObjectType是重复的 自己去钉钉申请
//			msg := TempAuthCodeCallbackMsg{AuthCorpId: "ding8ac2bab2b708b3cc35c2f4657eb6378f"}
//
//			orgId, dingtalkUserRoleBos, roleInitResp, teamId, err1 := initRelatedInfo(&msg, "test")
//
//			if err1 != nil {
//
//				fmt.Println("组装参数失败...............")
//				return
//			}
//
//			rootUserId, err1 := InitUserRoles(dingtalkUserRoleBos, orgId, teamId, &msg, roleInitResp)
//
//			convey.So(rootUserId, convey.ShouldNotBeNil)
//
//		})
//
//		if err != nil {
//			tx.Rollback()
//			fmt.Println(err)
//		}
//		tx.Commit()
//	}))
//}

//func TestTeamOwnerInit(t *testing.T) {
//
//	convey.Convey("Test TeamOwnerInit", t, test.StartUp(func(ctx context.Context) {
//
//		//开启事务
//		conn, err := mysql.GetConnect()
//		if err != nil {
//			log.Errorf(consts.DBOpenErrorSentence, err)
//		}
//		tx, err := conn.NewTx(nil)
//		if err != nil {
//			log.Errorf(consts.TxOpenErrorSentence, err)
//		}
//		defer mysql.Close(conn, tx)
//
//		convey.Convey("mock InitUserRoles teamId", func() {
//
//			//
//			////设置团队的owner
//			//err = inits.TeamOwnerInit(teamId, rootUserId, tx)
//			//
//			//
//			//convey.Convey("mock InitUserRoles",  tests.StartUpWithContext(func(ctx context.Context) {
//			//
//			//	msg:=TempAuthCodeCallbackMsg{AuthCorpId:"ding8ac2bab2b708b3cc35c2f4657eb6378f"}
//			//
//			//	orgId, dingtalkUserRoleBos, roleInitResp, teamId, err1 := initRelatedInfo(tx, &msg, "test")
//			//
//			//	if err1 !=nil {
//			//
//			//		fmt.Println("组装参数失败...............")
//			//		return
//			//	}
//			//
//			//	rootUserId, err1 := InitUserRoles(tx, dingtalkUserRoleBos, orgId, teamId, &msg, roleInitResp)
//			//
//			//	convey.So(rootUserId,convey.ShouldNotBeNil)
//			//
//			//}))
//
//			id := int64(1)
//
//			vo := orgvo.TeamOwnerInitReqVo{
//				TeamId: id,
//				Owner:  id}
//
//			err2 := orgfacade.TeamOwnerInit(vo)
//
//			if err2.Failure() {
//				fmt.Println("更新teamOwnerInit 失败 ..................")
//			}
//
//		})
//
//		if err != nil {
//			tx.Rollback()
//			fmt.Println(err)
//		}
//		tx.Commit()
//	}))
//
//}

func TestTeamUserInit(t *testing.T) {

	convey.Convey("Test TeamOwnerInit", t, test.StartUp(func(ctx context.Context) {

		//开启事务
		conn, err := mysql.GetConnect()
		if err != nil {
			log.Errorf(consts.DBOpenErrorSentence, err)
		}
		tx, err := conn.NewTx(nil)
		if err != nil {
			log.Errorf(consts.TxOpenErrorSentence, err)
		}
		defer mysql.Close(conn, tx)

		//inits.TeamUserInit(orgId, teamId, userId, isRoot, tx)

		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
		tx.Commit()
	}))
}

//func TestConstOrganizationCacheKeyFunc(t *testing.T) {
//	convey.Convey("Test ConstOrganizationCacheKeyFunc", t, test.StartUp(func(ctx context.Context) {
//
//		fmt.Println("原本的字符串:", consts.CacheUserConfig)
//
//		fmt.Printf("%q\n", strings.SplitAfter(consts.CacheUserConfig, consts.CacheKeyPrefix))
//
//		selfMap := make(map[string]interface{})
//
//		selfMap["orgId"] = 17
//		selfMap["userId"] = 1070
//
//		str, _ := util.ParseCacheKey(consts.CacheUserConfig, selfMap)
//
//		fmt.Println("结果:", str)
//
//	}))
//}
