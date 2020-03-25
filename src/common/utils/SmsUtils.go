package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	smsUtilOnce              sync.Once
	SmsAccessAppId           = ""
	SmsAccessSecret          = ""
	SmsSignName              = ""
	SmsTemplates             map[string]string //alicloud sms template map
	SmsSecurityLogTemplates  map[string]string //alicloud sms security log template map
	SmsSecurityLogsTemplates map[string]string //alicloud sms security logs template map
)

const (
	SmsTemplateCaptcha = iota
	SmsTemplateSecurityLog
	SmsTemplateSecurityLogs
)

func SmsInit() {
	smsUtilOnce.Do(func() {
		SmsAccessAppId = beego.AppConfig.String("sms::accessappid")
		SmsAccessSecret = beego.AppConfig.String("sms::accesssecret")
		SmsSignName = beego.AppConfig.String("sms::signname")

		SmsTemplates = make(map[string]string)
		SmsSecurityLogTemplates = make(map[string]string)
		SmsSecurityLogsTemplates = make(map[string]string)

		keys := []string{"templates", "templatesSecurityLog", "templatesSecurityLogs"}
		values := []map[string]string{SmsTemplates, SmsSecurityLogTemplates, SmsSecurityLogsTemplates}

		type templateItem struct {
			Lang string `json:"lang"`
			Name string `json:"name"`
		}
		type templateConfig struct {
			Items []*templateItem `json:"items"`
		}

		for i := 0; i < len(keys); i++ {
			jsonStr := beego.AppConfig.String("sms::" + keys[i])
			tmpConfig := &templateConfig{}
			if err := JsonParse(jsonStr, &tmpConfig); err != nil {
				logs.Error("json.Unmarshal sms::%s [%s] failed : %v", keys[i], jsonStr, err)
			}

			for _, item := range tmpConfig.Items {
				values[i][item.Lang] = item.Name
			}
		}
	})
}

type SendSmsReply struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

// SendSms 发送短信
func SendSms(phoneNumbers, templateParam, templateCode string, canReTry bool) error {

	paras := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", rand.Int63()),
		"AccessKeyId":      SmsAccessAppId,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",

		"Action":        "SendSms",
		"Version":       "2017-05-25",
		"RegionId":      "cn-hangzhou",
		"PhoneNumbers":  phoneNumbers,
		"SignName":      SmsSignName,
		"TemplateParam": templateParam,
		"TemplateCode":  templateCode,
	}

	var keys []string

	for k := range paras {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sortQueryString string

	for _, v := range keys {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, replace(v), replace(paras[v]))
	}

	stringToSign := fmt.Sprintf("GET&%s&%s", replace("/"), replace(sortQueryString[1:]))

	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", SmsAccessSecret)))
	mac.Write([]byte(stringToSign))
	sign := replace(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	str := fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s%s", sign, sortQueryString)

	resp, err := http.Get(str)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ssr := &SendSmsReply{}

	if err := json.Unmarshal(body, ssr); err != nil {
		return err
	}

	if ssr.Code == "SignatureNonceUsed" && canReTry {
		return SendSms(phoneNumbers, templateParam, templateCode, false)
	} else if ssr.Code != "OK" {
		return errors.New(ssr.Code)
	}

	return nil
}

func SmsMobileAndTemplate(tmpType int, nationalCode, mobile string) (mobileTmp, templateCode string, ok bool) {
	var smsLang string
	if nationalCode == "86" {
		mobileTmp = mobile //if it is mainland phone, national_code is not needed.
		smsLang = "zh"
	} else {
		mobileTmp = nationalCode + mobile
		smsLang = "en"
	}

	templateCode, ok = SmsTemplate(tmpType, smsLang)

	return
}

func SmsTemplate(tmpType int, lang string) (templateCode string, ok bool) {

	switch tmpType {
	case SmsTemplateCaptcha:
		{
			templateCode, ok = SmsTemplates[lang]
		}
	case SmsTemplateSecurityLog:
		{
			templateCode, ok = SmsSecurityLogTemplates[lang]
		}
	case SmsTemplateSecurityLogs:
		{
			templateCode, ok = SmsSecurityLogsTemplates[lang]
		}
	}

	if !ok {
		logs.Warning("no sms template found for lang [%s]", lang)
		return
	}

	return
}

func replace(in string) string {
	rep := strings.NewReplacer("+", "%20", "*", "%2A", "~", "%7E")
	return rep.Replace(url.QueryEscape(in))
}
