package date

import (
	"time"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/types"
)

func Format(t time.Time) string {
	return t.Format(consts.AppTimeFormat)
}
func FormatTime(t types.Time) string {
	return time.Time(t).Format(consts.AppTimeFormat)
}

func ParseTime(str string) types.Time {
	t, err := time.Parse(consts.AppTimeFormat, str)
	if err != nil {
		panic(err)
	}
	return types.Time(t)
}

func Parse(str string) time.Time {
	t, err := time.Parse(consts.AppTimeFormat, str)
	if err != nil {
		panic(err)
	}
	return t
}

// 获取某一天当前周的周一0点
func GetWeekStart(d time.Time) time.Time {
	offset := int(time.Monday - d.Weekday())
	if offset > 0 {
		offset = -6
	}

	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetMonthStart(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}
