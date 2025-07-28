/*
*
2 * @Author: Nico
3 * @Date: 2020/1/31 11:17
4
*/
package api

import (
	"context"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (r *mutationResolver) UserRegister(ctx context.Context, input vo.UserRegisterReq) (*vo.UserRegisterResp, error) {
	resp := orgfacade.UserRegister(orgvo.UserRegisterReqVo{
		Input: input,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Data, nil
}
