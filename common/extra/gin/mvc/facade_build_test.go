package mvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/model/vo/idvo"
)

type PostGreeter struct {
	Greeter
}

func (PostGreeter) ApplyPrimaryId(ctx *context.Context, req idvo.ApplyPrimaryIdReqVo) idvo.ApplyPrimaryIdRespVo {
	return idvo.ApplyPrimaryIdRespVo{}
}

func TestFacadeBuilder_Build(t *testing.T) {

	postGreeter := PostGreeter{Greeter: NewPostGreeter("idsvc", "127.0.0.1", 8080, "v1")}

	facadeBuilder := FacadeBuilder{
		StorageDir: "F:\\workspace-test",
		Package:    "facade",
		VoPackage:  "idvo",
		Greeters:   []interface{}{&postGreeter},
	}

	facadeBuilder.Build()
}
