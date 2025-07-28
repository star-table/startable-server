package pinyin

import (
	"testing"

	"github.com/mozillazg/go-pinyin"
)

func TestConvertCode(t *testing.T) {

	str := "这是One IssueOfFirst任务"
	t.Log(ConvertCode(str))

	str = "One这是OfIssue First1123456"
	t.Log(ConvertCode(str))

	str = "OneOneTwoTwoThirdThirdFourFourFive"
	t.Log(ConvertCode(str))
	t.Log(ConvertCodeWithMaxLen(str, 8))
}

func TestConvertToPinyin(t *testing.T) {

	strs := pinyin.Pinyin("刘s2千源", args)
	for k, str := range strs {
		t.Log(k, str)
	}
	t.Log(ConvertToPinyin("s23千源"))
	t.Log(ConvertToPinyin("abc"))
	t.Log(ConvertToPinyin("s23千源"))
	t.Log(ConvertToPinyin("123刘3123"))
	t.Log(ConvertToPinyin("1233123刘"))

}
