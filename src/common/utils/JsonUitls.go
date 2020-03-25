package utils

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

func JsonString(container interface{}) string {
	if container == nil {
		return "{}"
	}
	bytes, e := json.Marshal(container)
	if e != nil {
		beego.Warn(e.Error())
	}
	return string(bytes)
}
