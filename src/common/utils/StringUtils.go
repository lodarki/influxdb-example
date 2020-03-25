package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 获取随机n位数字
func RandomNum(len int) int64 {
	if len < 0 {
		return 0
	}

	// 对地位
	min := int64(math.Pow10(len - 1))
	max := int64(math.Pow10(len))

	rand.Seed(time.Now().UnixNano())
	result := rand.Int63n(max)

	if result <= min {
		result += min
	}

	return result
}

// 随机n位数字的字符串形式
func RandomNumStr(len int) string {
	randomInt64 := RandomNum(len)
	return strconv.FormatInt(randomInt64, 10)
}

// 对字符串md5加密
func GetMd5String(str string, salt ...string) string {
	t := md5.New()
	t.Write([]byte(str))
	for _, s := range salt {
		t.Write([]byte(s))
	}
	return hex.EncodeToString(t.Sum(nil))
}

func RandomNonceStr() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func IntToInt64(i int) int64 {
	itoa := strconv.Itoa(i)
	i64, err := strconv.ParseInt(itoa, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	return i64
}

func Int64ToInt(i int64) int {
	formatInt := strconv.FormatInt(i, 10)
	atoi, e := strconv.Atoi(formatInt)
	if e != nil {
		beego.Error(e)
	}
	return atoi
}

func JsonParse(jsonStr string, container interface{}) error {
	return json.Unmarshal([]byte(jsonStr), container)
}

func JsonParseMap(jsonStr string) (resultMap map[string]interface{}, err error) {
	resultMap = make(map[string]interface{})
	err = JsonParse(jsonStr, &resultMap)
	return
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func StringToInt64(str string) int64 {

	if len(str) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	return i
}

func StringToUint64(str string) uint64 {

	if len(str) == 0 {
		return 0
	}

	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	return i
}

func StringToUint32(str string) uint32 {
	return uint32(StringToUint64(str))
}

func StringToFloat64(str string) float64 {
	f, e := strconv.ParseFloat(str, 64)
	if e != nil {
		beego.Error(e)
	}
	return f
}

func StringToFloat32(str string) float64 {
	f, e := strconv.ParseFloat(str, 32)
	if e != nil {
		beego.Error(e)
	}
	return f
}

func SubTailStr(str string, index int) string {
	if len(str) > index {
		return str[(len(str) - index):]
	}
	return ""
}

// 可完整截取汉字，每一个汉字算1位
func SubRunesStr(str string, from, to int) string {

	runes := []rune(str)
	rlen := len(runes)

	if from > to {
		return ""
	}

	if from > rlen {
		return ""
	}

	if to > rlen {
		to = rlen
	}

	return string(runes[from:to])
}

func IsValidMobile(str string) bool {
	if len(str) == 0 {
		return false
	}

	regular := `^1[3456789]\d{9}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

// zp 需要补零到多少位
func IntToString(i int, zp int) string {
	itoa := strconv.Itoa(i)
	if len(itoa) >= zp {
		return itoa
	}

	for j := 0; j <= zp-len(itoa); j++ {
		itoa = "0" + itoa
	}
	return itoa
}

// 判断是否全是数字
func IsNumeric(val interface{}, valLen int) bool {

	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" || len(str) != valLen {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}
