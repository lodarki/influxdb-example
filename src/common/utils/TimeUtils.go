package utils

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strconv"
	"time"
)

const BaseFormatDateTimeM = "2006/01/02 15:04:05.000"
const BaseFormatDateTime = "2006-01-02 15:04:05"
const BaseFormatDate = "2006-01-02"
const BaseFormatMonth = "2006-01"
const BaseFormatTime = "15:04:05"
const BaseFormatHour = "15"
const DisplayFormatDateTime = "2006年01月02日 15点04分"
const DisplayFormatDate = "2006年01月02日"

var weekMap map[int]string

func init() {
	weekMap = make(map[int]string)
	weekMap[0] = "周日"
	weekMap[1] = "周一"
	weekMap[2] = "周二"
	weekMap[3] = "周三"
	weekMap[4] = "周四"
	weekMap[5] = "周五"
	weekMap[6] = "周六"
}

func GetWeekDayStr(date time.Time) string {
	return weekMap[int(date.Weekday())]
}

// 获取制定时间的零点
func GetZeroClockForDate(t time.Time) (time.Time, error) {
	return time.ParseInLocation(BaseFormatDate, t.Format(BaseFormatDate), time.Local)
}

func ParseMonth(dateStr string) (time.Time, error) {
	return time.ParseInLocation(BaseFormatMonth, dateStr, time.Local)
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.ParseInLocation(BaseFormatDate, dateStr, time.Local)
}

func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.ParseInLocation(BaseFormatDateTime, dateTimeStr, time.Local)
}

func ParseDateTimeM(dateTimeStr string) (time.Time, error) {
	return time.ParseInLocation(BaseFormatDateTimeM, dateTimeStr, time.Local)
}

// 判断时间是否比当前晚
func AfterNow(targetTime time.Time) bool {
	sub := time.Now().Sub(targetTime)
	return sub.Nanoseconds() < 0
}

// 累加小时
func AddHour(t time.Time, hour int) time.Time {
	return t.Add(time.Duration(hour) * time.Hour)
}

func MonthStr(t time.Time) string {
	return t.Format(BaseFormatMonth)
}

func DateStr(t time.Time) string {
	return t.Format(BaseFormatDate)
}

func DisplayDateStr(t time.Time) string {
	return t.Format(DisplayFormatDate)
}

func DisplayDateTimeStr(t time.Time) string {
	return t.Format(DisplayFormatDateTime)
}

func TimesStr(t time.Time) string {
	return t.Format(BaseFormatTime)
}

func DateTimeStr(t time.Time) string {
	return t.Format(BaseFormatDateTime)
}

func DateTimeMStr(t time.Time) string {
	return t.Format(BaseFormatDateTimeM)
}

// 解析时间
func HourStr(t time.Time) int {
	hour, _ := strconv.Atoi(t.Format(BaseFormatHour))
	return hour
}

// 解析时间
func TimeStr(t time.Time) int {
	time, _ := strconv.Atoi(t.Format(BaseFormatTime))
	return time
}

// Now 返回当前时间戳
func Now() int64 {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	cur := time.Now().In(loc).Unix()
	return cur
}

// UploadFilePath 上传文件路径
func UploadFilePath(exts ...string) string {
	date := time.Now()
	now := Now()
	var ext = "jpg"
	if len(exts) > 0 {
		ext = exts[0]
	}

	path := fmt.Sprintf("/%d/%02d/%d/%d."+ext, date.Year(), date.Month(), date.Day(), now)
	return path
}

func AddDay(day int) string {
	dateString := time.Now().AddDate(0, 0, day).String()
	ret, e := ParseDate(dateString)
	if e != nil {
		return ""
	}

	return ret.String()
}

// string to datetime
func StringToDateTimeHook(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	return data, nil
}

// 解析
func MapStructureDecode(target interface{}, sourceMap interface{}) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: StringToDateTimeHook,
		Result:     &target,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	err = decoder.Decode(sourceMap)
	if err != nil {
		return err
	}
	return nil
}

func GetFirstDayTimeOfMonth(targetTime time.Time) (t time.Time, err error) {
	t, err = ParseMonth(MonthStr(targetTime))
	return
}

func GetLastDayTimeOfMonth(targetTime time.Time) (t time.Time, err error) {
	firstDay, err := GetFirstDayTimeOfMonth(targetTime)
	if err != nil {
		return
	}
	t = firstDay.AddDate(0, 1, -1)
	return
}

func SubMonths(fromTime, toTime time.Time) int {
	return (toTime.Year()-fromTime.Year())*12 + int(toTime.Month()-fromTime.Month())
}

func SecondsOfTodayByMinAndSec(c, m int) (sec int64, err error) {

	restDay := c / 24
	restC := c % 24

	now := time.Now()
	timeStr := fmt.Sprintf("%v %v:%v:00", DateStr(now), IntToString(restC, 2), IntToString(m, 2))
	dateTime, e := ParseDateTime(timeStr)
	if e != nil {
		err = e
		return
	}

	if restDay > 0 {
		dateTime = dateTime.AddDate(0, 0, restDay)
	}
	sec = dateTime.Unix()
	return
}

func GetBlockTime(now time.Time) (target time.Time) {
	today, _ := GetZeroClockForDate(now)
	target = today.Add(time.Duration(now.Hour()) * time.Hour)
	m := now.Minute()

	if m < 15 {
		return
	}

	if m < 45 {
		target = target.Add(time.Duration(30) * time.Minute)
		return
	}

	target = target.Add(time.Duration(1) * time.Hour)
	return
}
