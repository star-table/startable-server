/*
*
2 * @Author: Nico
3 * @Date: 2020/1/31 11:17
4
*/
package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) UserRegister(req orgvo.UserRegisterReqVo) orgvo.UserRegisterRespVo {
	res, err := service.UserRegister(req)
	return orgvo.UserRegisterRespVo{Data: res, Err: vo.NewErr(err)}
}
