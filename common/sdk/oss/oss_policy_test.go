package oss

import (
	"testing"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/util/json"
)

func TestPostPolicy(t *testing.T) {

	config.LoadUnitTestConfig()
	pp := PostPolicy("project", 1000*60*5, 0)
	t.Log(json.ToJson(pp))
}

func TestPostPolicyWithCallback(t *testing.T) {
	config.LoadUnitTestConfig()
	pp := PostPolicyWithCallback("project", 1000*60*5, 0, "")
	t.Log(json.ToJson(pp))
}
