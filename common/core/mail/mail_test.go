package mail

import (
	"testing"

	"github.com/star-table/startable-server/common/core/config"
)

func TestAvailability(t *testing.T) {
	t.Logf("start load config")
	config.LoadUnitTestConfig()

	err := SendMail([]string{"ainililia@163.com"}, "hello", "hello")
	if err == nil {
		t.Log("successful")
	} else {
		t.Log(err)
		t.Log("send failed")
	}
}
