package util

import (
	"testing"

	"github.com/star-table/startable-server/common/core/util/slice"
)

func TestGetCurrentPath(t *testing.T) {

	t.Log(GetCurrentPath())

}

func TestSliceContain(t *testing.T) {
	list := []int64{1, 2, 3, 4}
	b, e := slice.Contain(list, 1)
	t.Log(b, e)
}
