package domain

import (
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// 删除钉钉酷应用
func DeleteDingCoolApp(orgId, projectId int64) {
	asyn.Execute(func() {
		orgfacade.DeleteCoolAppByProject(orgvo.DeleteCoolAppByProjectReq{
			OrgId:     orgId,
			ProjectId: projectId,
		})
	})
}

func UpdateDingTopCard(orgId, projectId int64) {
	if projectId <= 0 {
		return
	}
	asyn.Execute(func() {
		orgfacade.UpdateCoolAppTopCard(orgvo.UpdateCoolAppTopCardReq{
			OrgId:     orgId,
			ProjectId: projectId,
		})
	})
}
