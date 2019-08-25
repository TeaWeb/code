package teadb

import "github.com/TeaWeb/code/tealogs"

type AccessLogDAO interface {
	Init()

	TableName(day string) string
	FindAccessLogCookie(day string, logId string) (*tealogs.AccessLog, error)
	FindRequestHeaderAndBody(day string, logId string) (*tealogs.AccessLog, error)
	FindResponseHeaderAndBody(day string, logId string) (*tealogs.AccessLog, error)
	ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*tealogs.AccessLog, error)
	HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error)
	ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*tealogs.AccessLog, error)
	ListTopAccessLogs(day string, size int) ([]*tealogs.AccessLog, error)
	QueryAccessLogs(day string, serverId string, query *Query) ([]*tealogs.AccessLog, error)
}
