package callsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

//func TestPushSyncMsgToMq(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		t.Log(PushSyncMsgToMq())
//	}))
//}

func TestPushSyncMsgToMqHandlerFunc(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		key := consts.CacheFeiShuStopSyncScope
		isExist, isExistEr := cache.Exist(key)
		t.Log(isExistEr)
		t.Log(isExist)
	}))
}

func TestPushBotHelpNoticeForGroup(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		PushWelcomeMsgToOrgMember("1279794b670f575f", "ou_4736aa27afbeb880894f5f98e19ac46e")
	}))
}
