package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
)

func TestDepartments(t *testing.T) {
	convey.Convey("Test 获取部门列表", t, test.StartUp(func(ctx context.Context) {

		user, err := orgsvcService.GetCurrentUser(ctx)
		if err != nil {
			t.Fatalf("GetCurrentUser failed: %v", err)
		}

		convey.Convey("Test 获取部门列表", func() {
			isTop := 1

			reqVo := orgvo.DepartmentsReqVo{
				Params: &vo.DepartmentListReq{
					IsTop: &isTop,
				},
				OrgId:         user.OrgId,
				CurrentUserId: user.UserId,
			}

			// 创建测试请求
			reqBody, _ := json.Marshal(reqVo)
			req := httptest.NewRequest("POST", "/departments", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 创建Gin上下文
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 调用处理函数
			Departments(c)

			// 检查响应
			convey.So(w.Code, convey.ShouldEqual, http.StatusOK)
			t.Log(w.Body.String())
		})
	}))
}

func TestDepartmentMembers(t *testing.T) {
	convey.Convey("Test 获取部门成员列表", t, test.StartUp(func(ctx context.Context) {

		departmentId := int64(6)

		user, err := orgsvcService.GetCurrentUser(ctx)
		if err != nil {
			t.Fatalf("GetCurrentUser failed: %v", err)
		}

		convey.Convey("Test 获取部门成员列表", func() {
			reqVo := orgvo.DepartmentMembersReqVo{
				Params: vo.DepartmentMemberListReq{
					DepartmentID: &departmentId,
				},
				OrgId:         user.OrgId,
				CurrentUserId: user.UserId,
			}

			// 创建测试请求
			reqBody, _ := json.Marshal(reqVo)
			req := httptest.NewRequest("POST", "/departmentMembers", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 创建Gin上下文
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 调用处理函数
			DepartmentMembers(c)

			// 检查响应
			convey.So(w.Code, convey.ShouldEqual, http.StatusOK)
			t.Log(w.Body.String())
		})
	}))
}
