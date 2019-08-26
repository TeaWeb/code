package teadb

import (
	"github.com/TeaWeb/code/tealogs/accesslogs"
)

type AccessLogDAOInterface interface {
	// 初始化
	Init()

	// 获取表名
	TableName(day string) string

	// 写入一条日志
	InsertOne(accessLog *accesslogs.AccessLog) error

	// 写入一组日志
	InsertAccessLogs(accessLogList []interface{}) error

	// 查找某条访问日志的cookie信息
	FindAccessLogCookie(day string, logId string) (*accesslogs.AccessLog, error)

	// 查找某条访问日志的请求信息
	FindRequestHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error)

	// 查找某条访问日志的响应信息
	FindResponseHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error)

	// 列出日志
	ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error)

	// 检查是否有下一条日志
	HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error)

	// 判断某个代理服务是否有日志
	HasAccessLog(day string, serverId string) (bool, error)

	// 列出最近的某些日志
	ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*accesslogs.AccessLog, error)

	// 列出某天的一些日志
	ListTopAccessLogs(day string, size int) ([]*accesslogs.AccessLog, error)

	// 根据查询条件来查找日志
	QueryAccessLogs(day string, serverId string, query *Query) ([]*accesslogs.AccessLog, error)
}
