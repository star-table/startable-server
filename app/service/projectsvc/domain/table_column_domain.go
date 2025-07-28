package domain

import (
	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetTableColumns(orgId int64, userId int64, tableId int64, isNeedDescription bool) ([]*projectvo.TableColumnData, errs.SystemErrorInfo) {
	if tableId == 0 {
		// 找汇总表的表头
		summeryTableResp := tablefacade.GetSummeryTableId(projectvo.GetSummeryTableIdReqVo{
			OrgId:  orgId,
			UserId: 0,
			Input:  &tableV1.ReadSummeryTableIdRequest{},
		})
		if summeryTableResp.Failure() {
			log.Errorf("[GetTableColumns] failed, orgId: %d, tableId: %d, err: %v", orgId, tableId, summeryTableResp.Error())
			return nil, summeryTableResp.Error()
		}
		tableId = summeryTableResp.Data.TableId
	}
	tableColumnsResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:          []int64{tableId},
			IsNeedDescription: isNeedDescription,
		},
	})
	if tableColumnsResp.Failure() {
		log.Errorf("[GetTableColumns] err: %v", tableColumnsResp.Error())
		return nil, tableColumnsResp.Error()
	}

	if len(tableColumnsResp.Data.Tables) == 0 {
		return nil, errs.TableNotExist
	}
	return tableColumnsResp.Data.Tables[0].Columns, nil
}
