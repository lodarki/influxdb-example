package utils

import (
    "errors"
    "github.com/astaxie/beego"
    "reflect"
)

func SliceDel(arr interface{}, item interface{}) {
    val := reflect.ValueOf(arr)
    var value reflect.Value

    if val.Kind() == reflect.Ptr {
        value = reflect.Indirect(val)
    } else {
        value = val
    }

    if value.Kind() != reflect.Slice {
        beego.Error("invalid Param")
        return
    }

    var index = 0
    for index >= 0 {
        index = SliceIndexOf(arr, item)
        if index < 0 {
            return
        }

        oriLen := value.Len()
        if oriLen == 0 {
            return
        }
        var result []interface{}
        for i := 0; i < oriLen; i++ {
            if i != index {
                result = append(result, value.Index(i).Interface())
            }
        }

        value.SetLen(len(result))
        for i, v := range result {
            value.Index(i).Set(reflect.ValueOf(v))
        }
    }
}

func SliceIndexOf(arr interface{}, item interface{}) int {

    val := reflect.ValueOf(arr)
    var value reflect.Value

    if val.Kind() == reflect.Ptr {
        value = reflect.Indirect(val)
    } else {
        value = val
    }

    if value.Kind() != reflect.Slice {
        beego.Error("invalid Param")
        return -1
    }

    if value.Len() == 0 {
        return -1
    }

    for i := 0; i < value.Len(); i++ {
        if value.Index(i).Type() != reflect.TypeOf(item) {
            return -1
        }
        if value.Index(i).Interface() == reflect.ValueOf(item).Interface() {
            return i
        }
    }

    return -1
}

func SliceMax(arr interface{}) (res interface{}, err error) {
    val := reflect.ValueOf(arr)
    var value reflect.Value

    if val.Kind() == reflect.Ptr {
        value = reflect.Indirect(val)
    } else {
        value = val
    }

    if value.Kind() != reflect.Slice {
        err = errors.New("data must be slice")
        return
    }

    if value.Len() == 0 {
        err = errors.New("empty slice")
        return
    }

    resV := value.Index(0)
    for i := 0; i < value.Len(); i++ {

        subV := value.Index(i)
		subVK := subV.Kind()
		if subVK == reflect.String {
			if subV.String() > resV.String() {
				resV = subV
			}
        } else if SliceContains([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}, subVK) {
			if subV.Int() > resV.Int() {
				resV = subV
			}
		} else if SliceContains([]reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}, subVK) {
			if subV.Uint() > resV.Uint() {
				resV = subV
			}
		} else if SliceContains([]reflect.Kind{reflect.Float32, reflect.Float64}, subVK) {
			if subV.Float() > resV.Float() {
				resV = subV
			}
		}
    }

    res = resV.Interface()
    return
}

// 转置slice
func RevertSlice(arr interface{}) error {
    val := reflect.ValueOf(arr)
    if val.Kind() != reflect.Ptr {
        return errors.New("param must be a pointer, to make sure revaluing")
    }

    ind := reflect.Indirect(val)
    if ind.Kind() != reflect.Slice {
        return errors.New("invalid param")
    }

    var result = make([]interface{}, ind.Len())
    for i := 0; i < ind.Len(); i++ {
        result[ind.Len()-1-i] = ind.Index(i).Interface()
    }

    for j, v := range result {
        ind.Index(j).Set(reflect.ValueOf(v))
    }

    return nil
}

func SliceContains(arr interface{}, item interface{}) bool {
    return SliceIndexOf(arr, item) >= 0
}

func SliceContainArr(oriArr interface{}, desArr interface{}) bool {
	val := reflect.ValueOf(desArr)
	var value reflect.Value

	if val.Kind() == reflect.Ptr {
		value = reflect.Indirect(val)
	} else {
		value = val
	}

	if value.Kind() != reflect.Slice {
		return SliceIndexOf(oriArr, desArr) >= 0
	}

	if value.Len() == 0 {
		return true
	}

	var allIn = true
	for i := 0; i < value.Len(); i++ {
		if SliceIndexOf(oriArr, value.Index(i).Interface()) < 0 {
			allIn = false
			break
		}
	}

	return allIn
}

func SliceSerialInt(ori []int, step int, scope int) (result []int) {

    if len(ori) == 0 {
        return
    }

    var count = 1
    for i, v := range ori {
        if i > 0 && v <= ori[i-1]+step {
            count += 1
            result = append(result, v)
        } else {
            count = 1
            result = []int{v}
        }

        if count >= scope {
            return
        }
    }

    result = []int{}
    return
}

func SliceSerialInt64(ori []int64, step int64, scope int64) (result []int64) {

    if len(ori) == 0 {
        return
    }

    var count int64 = 1
    for i, v := range ori {
        if i > 0 && v <= ori[i-1]+step {
            count += 1
            result = append(result, v)
        } else {
            count = 1
            result = []int64{v}
        }

        if count >= scope {
            return
        }
    }

    result = []int64{}
    return
}

type OrderInt64Array []int64

func (o OrderInt64Array) Len() int {
    return len(o)
}

func (o OrderInt64Array) Less(i, j int) bool {
    return o[i] < o[j]
}

func (o OrderInt64Array) Swap(i, j int) {
    o[i], o[j] = o[j], o[i]
}

type OrderStringArray []string

func (o OrderStringArray) Len() int {
    return len(o)
}

func (o OrderStringArray) Less(i, j int) bool {
    return o[i] < o[j]
}

func (o OrderStringArray) Swap(i, j int) {
    o[i], o[j] = o[j], o[i]
}

// 间隔抽取数组中的元素，达到间隔过滤的效果。
func SliceSeparateFilter(array interface{}, count int) error {
    v := reflect.ValueOf(array)
    if v.Kind() != reflect.Ptr {
        return errors.New("array must be pointer")
    }

    ind := reflect.Indirect(v)
    if ind.Kind() != reflect.Slice {
        return errors.New("array must be slice")
    }

    if ind.Len() == 0 {
        return nil
    }

    if count <= 0 {
        return nil
    }

    if count >= ind.Len() {
        return nil
    }

    var result []interface{}
    var i = 0
    for i < ind.Len() {
        result = append(result, ind.Index(i).Interface())
        i = i + count
    }

    ind.SetLen(len(result))
    for i,v := range result {
        ind.Index(i).Set(reflect.ValueOf(v))
    }

    return nil
}