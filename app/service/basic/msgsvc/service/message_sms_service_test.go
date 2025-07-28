package msgsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/model/vo/msgvo"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
)

func TestSendEmail(t *testing.T) {
	convey.Convey("TestSendEmail", t, test.StartUp(func(ctx context.Context) {
		err := SendMail(msgvo.SendMailReqVo{
			Input: msgvo.SendMailReqData{
				Emails:  []string{"ainililia@163.com"},
				Subject: "text",
				Content: "body",
			},
		})
		assert.Equal(t, err, nil)
	}))

}
