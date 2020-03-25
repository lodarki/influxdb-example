package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"reflect"
	"strings"
	"time"
)

func ParseSlice(ori interface{}, des interface{}) (err error) {

	sliceVori := reflect.ValueOf(ori)

	if sliceVori.Kind() == reflect.Ptr {
		sliceVori = reflect.Indirect(sliceVori)
	}

	desv := reflect.ValueOf(des)
	if desv.Kind() != reflect.Ptr {
		return errors.New("des must be pointer of Slice ")
	}

	ind := reflect.Indirect(desv)
	if ind.Kind() != reflect.Slice {
		return errors.New("des must be pointer of Slice ")
	}

	for i := 0; i < sliceVori.Len(); i++ {
		vori := sliceVori.Index(i)
		sind := reflect.New(ind.Type().Elem())
		e := ParseStruct(vori.Interface(), sind.Interface())
		if e != nil {
			logs.Error(e)
		} else {
			ind = reflect.Append(ind, reflect.Indirect(sind))
		}
	}

	reflect.Indirect(desv).Set(ind)
	return nil
}

func ParseStruct(ori interface{}, des interface{}) (err error) {

	bytes, e := json.Marshal(ori)
	if e != nil {
		return e
	}

	var dataMap = make(map[string]interface{})
	e = json.Unmarshal(bytes, &dataMap)
	if e != nil {
		return e
	}

	if reflect.ValueOf(des).Kind() != reflect.Ptr {
		return errors.New("data must be a pointer")
	}

	ind := reflect.Indirect(reflect.ValueOf(des))

	switch ind.Kind() {
	case reflect.Map:
		return json.Unmarshal(bytes, des)
	case reflect.Struct:
		t := ind.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tagName := f.Tag.Get("json")
			d := dataMap[tagName]
			if d == nil {
				continue
			}

			switch ind.Field(i).Kind() {
			case reflect.String:
				var str = ""
				e := ForceConvert(&str, d)
				if e == nil {
					ind.Field(i).SetString(str)
				} else {
					beego.Warn(e)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var i64 int64
				e := ForceConvert(&i64, d)
				if e == nil {
					ind.Field(i).SetInt(i64)
				} else {
					beego.Warn(e)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var ui64 uint64
				e := ForceConvert(&ui64, d)
				if e == nil {
					ind.Field(i).SetUint(ui64)
				} else {
					beego.Warn(e)
				}
			case reflect.Float64, reflect.Float32:
				var f64 float64
				e := ForceConvert(&f64, d)
				if e == nil {
					ind.Field(i).SetFloat(f64)
				} else {
					beego.Warn(e)
				}
			default:
				if ind.Field(i).Type() == reflect.TypeOf(time.Time{}) {
					var t time.Time
					e := ForceConvert(&t, d)
					if e == nil {
						ind.Field(i).Set(reflect.ValueOf(t))
					} else {
						beego.Warn(e)
					}
				}
			}
		}
		break
	}

	return nil
}

// 强行转换赋值, 目前仅支持string，intX，uintX，floatX，time.Time 类型
func ForceConvert(container interface{}, ori interface{}) error {

	if ori == nil {
		return nil
	}

	v := reflect.ValueOf(container)
	if v.Kind() != reflect.Ptr {
		return errors.New("container must be pointer")
	}

	ind := reflect.Indirect(v)
	if reflect.TypeOf(ind) == reflect.TypeOf(ori) {
		ind.Set(reflect.ValueOf(ori))
		return nil
	}

	if ind.Kind() == reflect.String && reflect.TypeOf(ori) == reflect.TypeOf(time.Time{}) {
		ind.SetString(DateTimeStr(ori.(time.Time)))
		return nil
	} else if reflect.ValueOf(ori).Kind() == reflect.Struct {
		return errors.New("unmatched struct")
	}

	oriStr := fmt.Sprintf("%v", ori)
	oriV := reflect.ValueOf(ori)
	switch oriV.Kind() {
	case reflect.Float32, reflect.Float64:
		i := int64(oriV.Float())
		oriStr = fmt.Sprintf("%v", i)
		break
	}

	switch ind.Kind() {
	case reflect.String:
		ind.SetString(oriStr)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ind.SetInt(StringToInt64(oriStr))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ind.SetUint(StringToUint64(oriStr))
		return nil
	case reflect.Float32, reflect.Float64:
		ind.SetFloat(StringToFloat64(oriStr))
		return nil
	}

	if ind.Type() == reflect.TypeOf(time.Time{}) {
		index := strings.Index(oriStr, ".")
		if index > 0 {
			oriStr = oriStr[:index]
		}
		t, e := ParseDateTime(oriStr)
		if e == nil {
			ind.Set(reflect.ValueOf(t))
			return nil
		} else {
			beego.Error(e)
		}
	}

	return errors.New("un support type")
}
