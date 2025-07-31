package orgsvc

import (
	"fmt"

	"github.com/google/martian/log"
	"github.com/spf13/cast"
	"github.com/star-table/startable-server/app/service/orgsvc/po"
	"github.com/star-table/startable-server/common/core/consts"

	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/model/vo/commonvo"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InitOrg(initOrgBo bo.InitOrgBo) (int64, errs.SystemErrorInfo) {
	cacheKey := fmt.Sprintf("%s%s", sconsts.CacheFsOrgInit, initOrgBo.OutOrgId)
	defer cache.Del(cacheKey)
	cache.SetEx(cacheKey, initOrgBo.OutOrgId, 60)

	var orgInfo *po.PpmOrgOrganization
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		var err errs.SystemErrorInfo
		orgInfo, err = domain.InitOrg(initOrgBo, tx)
		if err != nil {
			return err
		}
		orgRemarkJson := orgInfo.Remark
		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
		if len(orgRemarkJson) > 0 {
			oriErr := json.FromJson(orgRemarkJson, orgRemarkObj)
			if oriErr != nil {
				return errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
			}
		}
		// 将汇总表的 appId 存入组织属性（remark）中。如果汇总表等应用不存在，则创建之。
		_, _, err = SaveOrgSomeTableAppId(orgInfo.Id, 0, orgRemarkObj, tx)
		if err != nil {
			log.Error(err)
			return err
		}

		return err
	})
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.OrgInitError, err)
	}
	// 以下的操作不是很重要，不需要放到事务里面？这是之前的逻辑，先放到外面吧
	// 创建应用接口暂不支持创建文件夹，因此默认 0
	userId := int64(0)
	// 同步无码组织配置和组织字段
	saveErr := saveOrgFields(orgInfo.Id, userId)
	if saveErr != nil {
		log.Error(saveErr)
		return 0, saveErr
	}

	// 创建组织时，创建项目视图及其视图镜像
	// 极星标品：创建目录应用、创建任务视图
	//if viewResp := projectfacade.SyncOrgDefaultViewMirror(projectvo.SyncOrgDefaultViewMirrorReq{
	//	OrgIds:                              []int64{orgInfo.Id},
	//	StartPage:                           1,
	//	NeedUpdateSummaryAppVisibilityToAll: true,
	//}); viewResp.Failure() {
	//	log.Error(viewResp.Error())
	//	return 0, viewResp.Error()
	//}
	viewsResp := projectfacade.CreateOrgDirectoryAppsAndViews(projectvo.CreateOrgDirectoryAppsReq{OrgId: orgInfo.Id})
	if viewsResp.Failure() {
		log.Errorf("[InitOrg] projectfacade.CreateOrgDirectoryAppsAndViews err:%v, orgId:%v",
			viewsResp.Error(), orgInfo.Id)
		return 0, viewsResp.Error()
	}

	// 事件上报
	asyn.Execute(func() {
		e := &commonvo.OrgEvent{}
		e.OrgId = orgInfo.Id
		e.New = orgInfo

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		report.ReportOrgEvent(msgPb.EventType_OrgInited, openTraceIdStr, e)
	})

	return orgInfo.Id, nil
}

func SendFeishuMemberHelpMsg(tenantKey string, orgId int64, ownerOpenId string) errs.SystemErrorInfo {
	page := 1
	size := 50
	for {
		bos, _, err := domain.GetOrgUserInfoListBySourceChannel(orgId, sdk_const.SourceChannelFeishu, page, size)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(bos) == 0 {
			break
		}

		needIds := []string{}
		for _, info := range bos {
			//管理员不需要推送欢迎语（已推送安装提示）
			if info.OutUserId != "" && info.OutUserId != ownerOpenId {
				needIds = append(needIds, info.OutUserId)
			}
		}
		if len(needIds) > 0 {
			domain.FeishuMemberHelpMsgToMq(bo.FeishuHelpObjectBo{OrgId: orgId, TenantKey: tenantKey, OpenIds: needIds})
		}

		if len(bos) < size {
			break
		} else {
			page++
		}
	}

	return nil
}
