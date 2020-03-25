package constants

import "sync"

const (
	ClientTypeUnkown = iota
	ClientTypeWeb
	ClientTypeAdminWeb
	ClientTypeMonitor
	ClientTypeCmd
)

var clientTypeOnce sync.Once
var ClientTypeMap map[int]string
var ClientTypeDescMap map[string]int

func init() {
	clientTypeOnce.Do(func() {

		ClientTypeMap = make(map[int]string)
		ClientTypeMap[ClientTypeUnkown] = "Unkown"
		ClientTypeMap[ClientTypeWeb] = "Web"
		ClientTypeMap[ClientTypeAdminWeb] = "admin_web"
		ClientTypeMap[ClientTypeMonitor] = "Monitor"
		ClientTypeMap[ClientTypeCmd] = "Cmd"

		ClientTypeDescMap = make(map[string]int)
		ClientTypeDescMap["Unkown"] = ClientTypeUnkown
		ClientTypeDescMap["Web"] = ClientTypeWeb
		ClientTypeDescMap["admin_web"] = ClientTypeAdminWeb
		ClientTypeDescMap["Monitor"] = ClientTypeMonitor
		ClientTypeDescMap["Cmd"] = ClientTypeCmd
	})
}
