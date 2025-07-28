package callsvc

import (
	"fmt"
	"net/http"

	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/consts"
	consts1 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

var log = logger.GetDefaultLogger()

func UploadCallbackHandlerFunc(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		UploadCallbackHandler(c.Writer, c.Request)
	})
}

func UploadCallbackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
		}
	}()

	if r.Method == "POST" {
		fmt.Println("\nHandle Post Request...")

		// Get PublicKey bytes
		bytePublicKey, err := getPublicKey(r)
		if err != nil {
			responseFailed(w, getErrorRespBody(err.Error()))
			return
		}

		// Get Authorization bytes : decode from Base64String
		byteAuthorization, err := getAuthorization(r)
		if err != nil {
			responseFailed(w, getErrorRespBody(err.Error()))
			return
		}

		// Get MD5 bytes from Newly Constructed Authrization String.
		body, byteMD5, err := getMD5FromNewAuthString(r)
		if err != nil {
			responseFailed(w, getErrorRespBody(err.Error()))
			return
		}

		// VerifySignature and response to client
		if verifySignature(bytePublicKey, byteMD5, byteAuthorization) {
			form := r.PostForm
			fmt.Println(json.ToJson(form))

			// Do something you want accoding to callback_body ...
			callbackBody := &bo.OssCallBackBody{}
			jsonErr := json.FromJson(body, callbackBody)
			if jsonErr != nil {
				log.Error(jsonErr)
				responseFailed(w, getErrorRespBody(jsonErr.Error()))
				return
			}
			log.Infof("[UploadCallbackHandler] orgId:%d, projectId:%d, policyType:%d", callbackBody.OrgId, callbackBody.ProjectId, callbackBody.Type)

			bucket := callbackBody.Bucket
			var resourceId int64

			switch callbackBody.Type {
			case consts1.OssPolicyTypeIssueResource, consts1.OssPolicyTypeLesscodeResource, consts1.OssPolicyTypeCommentAttachments:
				path := util.JointUrl(callbackBody.Host, callbackBody.Object)
				respVo := projectfacade.CreateIssueResource(projectvo.CreateIssueResourceReqVo{
					UserId: callbackBody.UserId,
					OrgId:  callbackBody.OrgId,
					Input: vo.CreateIssueResourceReq{
						// 项目id
						ProjectId: callbackBody.ProjectId,
						// 任务id
						IssueId: callbackBody.IssueId,
						// 资源路径
						ResourcePath: path,
						// 资源大小，单位B
						ResourceSize: callbackBody.Size,
						// 文件名
						FileName: callbackBody.RealName,
						// 文件后缀
						FileSuffix: callbackBody.Format,
						// bucketName
						BucketName: &bucket,
						// 来源类型
						PolicyType: callbackBody.Type,
						// 文件存储类型
						ResourceType: consts1.OssResource,
					},
				})
				if respVo.Failure() {
					log.Error(respVo.Message)
					responseFailed(w, getErrorRespBody(respVo.Message))
					return
				}
				resourceId = respVo.Void.ID

			case consts1.OssPolicyTypeProjectResource:
				respVo := projectfacade.CreateProjectResource(projectvo.CreateProjectResourceReqVo{
					UserId: callbackBody.UserId,
					OrgId:  callbackBody.OrgId,
					Input: vo.CreateProjectResourceReq{
						ProjectID: callbackBody.ProjectId,
						FolderID:  callbackBody.FolderId,
						// 资源路径
						ResourcePath: util.JointUrl(callbackBody.Host, callbackBody.Object),
						// 资源大小，单位B
						ResourceSize: callbackBody.Size,
						// 文件名
						FileName: callbackBody.RealName,
						// 文件后缀
						FileSuffix: callbackBody.Format,
						// bucketName
						BucketName: &bucket,
						// 来源类型
						PolicyType: callbackBody.Type,
						// 文件存储类型
						ResourceType: consts1.OssResource,
					},
				})
				if respVo.Failure() {
					log.Error(respVo.Message)
					responseFailed(w, getErrorRespBody(respVo.Message))
					return
				}
				resourceId = respVo.Void.ID
			}

			// response OK : 200
			responseSuccess(w, getSuccessRespBody(resourceId))
		} else {
			// response FAILED : 400
			responseFailed(w, getErrorRespBody("无效的oss回调签名"))
		}
	}
}

func getErrorRespBody(msg string) string {
	return "{\"errors\":[{\"message\":\"{\\\"code\\\":500,\\\"message\\\":\\\"" + msg + "\\\"}\"}]}"
}

func getSuccessRespBody(resourceId int64) string {
	return fmt.Sprintf("{\"data\":{\"ossCallbackStatus\":200, \"resourceId\":%d}}", resourceId)
}
