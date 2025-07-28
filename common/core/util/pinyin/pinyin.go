package pinyin

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

var args = pinyin.NewArgs()
var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")

func ConvertCode(name string) string {
	target := ""

	var needAppend = true
	var lastCh = 'a'
	for _, ch := range name {
		if ch == ' ' {
			needAppend = true
			continue
		}

		if matchPinyin(&target, &ch) {
			continue
		}

		if matchASCII(&target, &ch) {
			continue
		}

		if dealAppend(&target, &ch, &lastCh, &needAppend) {
			continue
		}

		lastCh = ch
	}

	return target
}

// 判断是否是汉字,转拼音
func matchPinyin(target *string, ch *int32) (isContinue bool) {

	if hzRegexp.MatchString(string(*ch)) {
		strs := pinyin.Pinyin(string(*ch), args)
		if len(strs) > 0 {
			if len(strs[0]) > 0 {
				py := strs[0][0]
				*target += strings.ToUpper(string(py[0]))
			}
		}
		isContinue = true
	}

	return isContinue
}

// 判断是否是
func matchASCII(target *string, ch *int32) (isContinue bool) {
	if *ch >= '0' && *ch <= '9' {
		*target += string(*ch)
		isContinue = true
	}
	return isContinue
}

// 处理append标识和字符串拼接
func dealAppend(target *string, ch *int32, lastCh *int32, needAppend *bool) (isContinue bool) {

	if (*ch >= 'a' && *ch <= 'z') || (*ch >= 'A' && *ch <= 'Z') {
		if *needAppend {
			*needAppend = false
			*target += strings.ToUpper(string(*ch))
		} else if ((*lastCh >= 'a' && *lastCh <= 'z') || (*lastCh >= 'A' && *lastCh <= 'Z') || (hzRegexp.MatchString(string(*lastCh)))) && *ch >= 'A' && *ch <= 'Z' {
			*target += strings.ToUpper(string(*ch))
		}
		isContinue = true
	}
	return isContinue
}

func ConvertCodeWithMaxLen(name string, maxLen int) string {
	target := ConvertCode(name)
	if len(target) > maxLen {
		return target[0:8]
	} else {
		return target
	}
}

//func ConvertCode(name string) string {
//	target := ""
//
//	var needAppend = true
//	var lastCh = 'a'
//	for _, ch := range name {
//		if ch == ' ' {
//			needAppend = true
//			continue
//		}
//
//		if hzRegexp.MatchString(string(ch)) {
//			strs := pinyin.Pinyin(string(ch), args)
//			if len(strs) > 0 {
//				if len(strs[0]) > 0 {
//					py := strs[0][0]
//					target += strings.ToUpper(string(py[0]))
//				}
//			}
//			continue
//		}
//
//		if ch >= '0' && ch <= '9' {
//			target += string(ch)
//			continue
//		}
//
//		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
//			if needAppend {
//				needAppend = false
//				target += strings.ToUpper(string(ch))
//			} else if ((lastCh >= 'a' && lastCh <= 'z') || (lastCh >= 'A' && lastCh <= 'Z') || (hzRegexp.MatchString(string(lastCh)))) && ch >= 'A' && ch <= 'Z' {
//				target += strings.ToUpper(string(ch))
//			}
//			continue
//		}
//
//		lastCh = ch
//	}
//
//	return target
//}

func ConvertToPinyin(str string) string {
	if str == "" {
		return str
	}
	result := ""
	for _, c := range str {
		target := string(c)
		strs := pinyin.Pinyin(target, args)
		if len(strs) != 0 {
			target = Capitalize(strs[0][0])
		}
		result += target
	}
	return result
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	if str == "" {
		return str
	}
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
