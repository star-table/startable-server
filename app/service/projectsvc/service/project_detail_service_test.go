package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/copyer"
	bo2 "github.com/star-table/startable-server/common/model/bo"
)

func TestUpdateProjectDetail(t *testing.T) {

	bo := &bo2.ProjectDetailBo{}
	bo.UpdateTime = time.Now()
	po := &po.PpmProProjectDetail{}
	mayBlank, _ := time.Parse(consts.AppTimeFormat, "")
	fmt.Println(mayBlank.Format(consts.AppTimeFormat))
	copyer.Copy(bo, po)

	t.Log(po)

}
