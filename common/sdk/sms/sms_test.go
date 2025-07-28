package sms

import (
	"testing"

	"github.com/star-table/startable-server/common/core/config"
)

func TestSendSMS(t *testing.T) {

	config.LoadUnitTestConfig()
	resp, err := SendSMS("13122323528", "北极星", "SMS_175533634", map[string]string{
		"code": "TEST",
	})
	t.Log(err)
	t.Log(resp)
	t.Log(resp.Message)
}
