package controllers

import (
	"common/constants"
	"common/entitys"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type BaseController struct {
	beego.Controller
	DefaultData interface{}
	o           orm.Ormer
	lang        string
}

// 获取session
func (c *BaseController) GetSession() session.Store {
	globalSession := beego.GlobalSessions
	sess, err := globalSession.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		logs.Error("err", err.Error())
	}
	return sess
}

func (c *BaseController) GetDateTime(key string, def ...time.Time) (t time.Time, err error) {
	tStr := c.GetString(key)
	if len(tStr) == 0 {
		return
	}

	t, err = utils.ParseDateTime(tStr)
	if err != nil && len(def) > 0 {
		t = def[0]
	}
	return
}

func (c *BaseController) GetDate(key string, def ...time.Time) (t time.Time, err error) {
	tStr := c.GetString(key)
	if len(tStr) == 0 {
		if len(def) > 0 {
			t = def[0]
			return
		} else {
			err = errors.New("empty date")
			return
		}
	}

	t, err = utils.ParseDate(tStr)
	if err != nil && len(def) > 0 {
		t = def[0]
	}
	return
}

func (c *BaseController) GetAuthorization() string {
	tokenStr := c.Ctx.Input.Header("Authorization")
	if tokenStr == "" {
		tokenStr = c.GetString("token")
	}

	return tokenStr
}

// 获取IP
func (c *BaseController) GetIp() string {
	return c.Ctx.Input.IP()
}

// @Description get client type
func (c *BaseController) clientType() (t int) {
	t = constants.ClientTypeUnkown
	s := c.Ctx.Input.Header("client-type")
	if len(s) <= 0 {
		//LOG_FUNC_DEBUG("clientType %s unkown.", s)
		s = constants.ClientTypeMap[constants.ClientTypeWeb]
	}

	if tp, ok := constants.ClientTypeDescMap[s]; ok {
		t = tp
	}

	if t == constants.ClientTypeUnkown {
		logs.Warn("clientType %s unkown.", s)
		return
	}

	return
}

// @Description get access_token
func (c *BaseController) accessToken() string {
	tokenStr := c.Ctx.Input.Header("token")
	if len(tokenStr) == 0 {
		tokenStr = c.Ctx.Input.Cookie("token")
	}
	return tokenStr
}

// 提交
func (c *BaseController) Commit() {
	err := c.o.Commit()
	logs.Error("commit", err)
}

// 回滚
func (c *BaseController) Rollback() {
	err := c.o.Rollback()
	logs.Error("rollback", err)
}

// 回调函数
func (c *BaseController) FailFn() {
}

// 成功回到函数
func (c *BaseController) SuccessFn() {
}

// OkResult 返回正确结果
func (c *BaseController) Success(data interface{}, encoding ...bool) {
	if data == nil {
		data = "{}"
	}
	Message := "success"
	result := entitys.ApiResult{Code: entitys.ERROR_CODE_SUCCESS, Result: data, Msg: Message}
	c.Data["json"] = &result

	c.SuccessFn()
	encode := true
	if len(encoding) > 0 {
		encode = encoding[0]
	}
	c.ServeJSON(encode)
}

// ErrResult 错误返回结果
func (c *BaseController) Fail(code entitys.ErrorCode, res interface{}, message string) {
	result := entitys.ApiResult{Code: code, Result: res, Msg: message}
	c.Data["json"] = &result
	c.FailFn()
	c.ServeJSON(true)
}

// 请求错误
func (c *BaseController) BadRequest(message string) {
	c.Fail(entitys.ERROR_CODE_INVALID_ARGUMENT, c.DefaultData, message)
}

// 服务器内部错误
func (c *BaseController) ServerError(message string) {
	c.Fail(entitys.ERROR_CODE_INTERNAL, c.DefaultData, message)
}

// 没有权限
func (c *BaseController) NotPermission(message string) {
	c.Fail(entitys.ERROR_CODE_PERMISSION_DENIED, c.DefaultData, message)
}

func (c *BaseController) NeedToLogin(message string) {
	c.Fail(entitys.ERROR_CODE_COMMON_UNLOGIN, c.DefaultData, message)
}

// 数据库执行错误
func (c *BaseController) SQLError(message string) {
	c.Fail(entitys.ERROR_CODE_DB, c.DefaultData, message)
}

// 查询无数据
func (c *BaseController) DbNoData(message string) {
	c.Fail(entitys.ERROR_CODE_NOT_DATA_FOUND, c.DefaultData, message)
}

// 获取当前URL
func (c *BaseController) CurrentUrl() string {
	uri := c.Ctx.Input.URI()
	return c.Ctx.Input.Site() + uri
}

// 获取当前域名
func (c *BaseController) Host() string {
	return c.Ctx.Input.Site()
}

// 解析json
func (c *BaseController) JsonDecodeRequestBody(data interface{}) error {
	return json.Unmarshal(c.Ctx.Input.RequestBody, data)
}

type JsonParam map[string]interface{}

func (c *BaseController) JsonParamMap() JsonParam {
	var paramMap = make(map[string]interface{})
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &paramMap)
	return paramMap
}

func (c *BaseController) DownloadFile(filePath string) error {

	file, e := os.Open(filePath)
	if e != nil {
		return e
	}

	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", path.Base(filePath))) // 文件名
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	// 最主要的一句
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
	_ = file.Close()
	return nil
}

func (c *BaseController) GetJsonParam(key string, d interface{}, jsonParams ...JsonParam) error {
	var jsonParam JsonParam
	if len(jsonParams) > 0 {
		jsonParam = jsonParams[0]
	} else {
		jsonParam = c.JsonParamMap()
	}

	bytes, marshalE := json.Marshal(jsonParam[key])
	if marshalE != nil {
		return marshalE
	}
	return json.Unmarshal(bytes, d)
}

func (c *BaseController) GetBindStrings(key string) (result []string, bindE error) {
	bindE = c.Ctx.Input.Bind(&result, key)
	if len(result) == 0 {
		keyJson := c.GetString(key)
		bindE = json.Unmarshal([]byte(keyJson), &result)
	}
	return
}

func (c *BaseController) GetBindInt64s(key string) (result []int64, bindE error) {
	bindE = c.Ctx.Input.Bind(&result, key)
	if len(result) == 0 {
		keyJson := c.GetString(key)
		bindE = json.Unmarshal([]byte(keyJson), &result)
	}
	return
}

func (c *BaseController) GetBindInts(key string) (result []int, bindE error) {
	bindE = c.Ctx.Input.Bind(&result, key)
	if len(result) == 0 {
		keyJson := c.GetString(key)
		bindE = json.Unmarshal([]byte(keyJson), &result)
	}
	return
}

func (c *BaseController) GetSearchKeyValue(enableKeys []string) (searchKey, searchValue string, valid bool) {

	valid = true
	searchKey = c.GetString("search_key")
	searchKey = strings.ToLower(searchKey)
	searchValue = c.GetString("search_value")

	if len(searchValue) == 0 {
		searchKey = ""
		return
	}

	if len(searchKey) > 0 {
		if !utils.SliceContains(enableKeys, searchKey) {
			c.BadRequest("invalid search_key")
			valid = false
			return
		}
	}

	return
}

func (c *BaseController) GetOrderType(enableOrderTypes []string) (orderType string, valid bool) {
	orderType = c.GetString("order_type")
	orderType = strings.ToLower(orderType)

	valid = true
	if len(orderType) == 0 {
		return
	}

	if !utils.SliceContains(enableOrderTypes, orderType) {
		c.BadRequest("invalid order_type")
		valid = false
		return
	}

	return
}

func (c *BaseController) CheckCaptcha(k, v string) bool {
	code, err := getCaptcha(k)
	if err != nil {
		return false
	}

	if code != v {
		return false
	}

	//if captcha is correct, remove it. make sure it is used only once.
	_ = removeCaptcha(k)

	return true
}

func getCaptcha(k string) (code string, err error) {
	var (
		v  interface{}
		ok bool
	)

	if v = utils.CacheForMemory.Get(k); v == nil { // 查找失败
		logs.Debug("getCaptcha of %s failed : %v", k)
		err = errors.New(fmt.Sprintf("getCaptcha of %s ", k))
		return
	} else if code, ok = v.(string); !ok { // 转换失败
		logs.Error("cast captcha [%v] to string failed.", v)
		err = errors.New(fmt.Sprintf("cast captcha [%v] to string failed.", v))
		return
	}

	return
}

func removeCaptcha(k string) (err error) {
	err = utils.CacheForMemory.Delete(k)
	return
}

// @Description get lang
func (c *BaseController) Lang() string {
	//parse only the first time
	if len(c.lang) <= 0 {
		if v := c.Ctx.Input.Header("lang"); v != "" {
			c.lang = v
		} else {
			c.lang = "en"
		}
	}

	return c.lang
}
