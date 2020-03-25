package constants

const (
	RoleIdAdmin        = 1 //超级管理员
	RoleIdManager      = 2 //高级网员
	RoleIdNormal       = 3 //监控员
	RoleIdLocalAdmin   = 4 //监控端本地超级管理员
	RoleIdLocalManager = 5 //监控端本地高级网员
	RoleIdLocalNormal  = 6 //监控端本地监控员
	RoleIdMonitor      = 7 //网管

	RolePowerView = 0 //查看权限

	MonitorRemoteAble   = 1 //监控端支持远程控制
	MonitorRemoteUnAble = 2 //监控端不支持远程控制
)

const (
	RoleStatusStop   = iota //角色禁用状态
	RoleStatusNormal        //角色正常状态
)
