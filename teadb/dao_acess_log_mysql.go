package teadb

import (
	"context"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/TeaWeb/code/tealogs/accesslogs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type MySQLAccessLogDAO struct {
}

// 初始化
func (this *MySQLAccessLogDAO) Init() {
	return
}

// 获取表名
func (this *MySQLAccessLogDAO) TableName(day string) string {
	this.initTable("logs." + day)
	return "logs." + day
}

// 获取当前时间表名
func (this *MySQLAccessLogDAO) TodayTableName() string {
	return this.TableName(timeutil.Format("Ymd"))
}

// 写入一条日志
func (this *MySQLAccessLogDAO) InsertOne(accessLog *accesslogs.AccessLog) error {
	if accessLog.Id.IsZero() {
		accessLog.Id = shared.NewObjectId()
	}
	return NewQuery(this.TodayTableName()).
		InsertOne(accessLog)
}

// 写入一组日志
func (this *MySQLAccessLogDAO) InsertAccessLogs(accessLogList []interface{}) error {
	return NewQuery(this.TodayTableName()).
		InsertOnes(accessLogList)
}

// 查找某条访问日志的cookie信息
func (this *MySQLAccessLogDAO) FindAccessLogCookie(day string, logId string) (*accesslogs.AccessLog, error) {
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", logId).
		Result("_id", "cookie").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

// 查找某条访问日志的请求信息
func (this *MySQLAccessLogDAO) FindRequestHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", logId).
		Result("_id", "header", "requestData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

// 查找某条访问日志的响应信息
func (this *MySQLAccessLogDAO) FindResponseHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", logId).
		Result("_id", "sentHeader", "responseBodyData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

// 列出日志
func (this *MySQLAccessLogDAO) ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		query.Lt("_id", fromId)
	}
	if onlyErrors {
		query.Or([]OperandMap{
			{
				"hasErrors": []*Operand{NewOperand(OperandEq, true)},
			},
			{
				"status": []*Operand{NewOperand(OperandGte, 400)},
			},
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}
	query.Offset(offset)
	query.Limit(size)
	query.Desc("_id")
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

// 检查是否有下一条日志
func (this *MySQLAccessLogDAO) HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId).
		Result("_id")
	if len(fromId) > 0 {
		query.Lt("_id", fromId)
	}
	if onlyErrors {
		query.Or([]OperandMap{
			{
				"hasErrors": []*Operand{NewOperand(OperandEq, true)},
			},
			{
				"status": []*Operand{NewOperand(OperandGte, 400)},
			},
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}

	one, err := query.FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return false, err
	}
	return one != nil, nil
}

// 判断某个代理服务是否有日志
func (this *MySQLAccessLogDAO) HasAccessLog(day string, serverId string) (bool, error) {
	query := NewQuery(this.TableName(day))
	one, err := query.Attr("serverId", serverId).
		Result("_id").
		FindOne(new(accesslogs.AccessLog))
	return one != nil, err
}

// 列出最近的某些日志
func (this *MySQLAccessLogDAO) ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))

	shouldReverse := true
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		query.Gt("_id", fromId)
		query.Asc("_id")
	} else {
		query.Desc("_id")
		shouldReverse = false
	}
	if onlyErrors {
		query.Or([]OperandMap{
			{
				"hasErrors": []*Operand{NewOperand(OperandEq, true)},
			},
			{
				"status": []*Operand{NewOperand(OperandGte, 400)},
			},
		})
	}
	query.Limit(size)
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	if shouldReverse {
		lists.Reverse(ones)
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}

	return result, nil
}

// 列出某天的一些日志
func (this *MySQLAccessLogDAO) ListTopAccessLogs(day string, size int) ([]*accesslogs.AccessLog, error) {
	ones, err := NewQuery(this.TableName(day)).
		Limit(size).
		Desc("_id").
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

// 根据查询条件来查找日志
func (this *MySQLAccessLogDAO) QueryAccessLogs(day string, serverId string, query *Query) ([]*accesslogs.AccessLog, error) {
	query.table = this.TableName(day)
	ones, err := query.
		Attr("serverId", serverId).
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MySQLAccessLogDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	conn, err := SharedDB().(*MySQLDriver).connect()
	if err != nil {
		return
	}

	_, err = conn.ExecContext(context.Background(), "SHOW CREATE TABLE `"+table+"`")
	if err != nil {
		s := "CREATE TABLE `" + table + "` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`_id` varchar(24) DEFAULT NULL," +
			"`serverId` varchar(64) DEFAULT NULL," +
			"`backendId` varchar(64) DEFAULT NULL," +
			"`locationId` varchar(64) DEFAULT NULL," +
			"`fastcgiId` varchar(64) DEFAULT NULL," +
			"`rewriteId` varchar(64) DEFAULT NULL," +
			"`teaVersion` varchar(32) DEFAULT NULL," +
			"`remoteAddr` varchar(64) DEFAULT NULL," +
			"`remotePort` int(11) unsigned DEFAULT '0'," +
			"`remoteUser` varchar(128) DEFAULT NULL," +
			"`requestURI` varchar(1024) DEFAULT NULL," +
			"`requestPath` varchar(1024) DEFAULT NULL," +
			"`requestLength` bigint(20) unsigned DEFAULT '0'," +
			"`requestTime` decimal(20,6) unsigned DEFAULT '0.000000'," +
			"`requestMethod` varchar(16) DEFAULT NULL," +
			"`requestFilename` varchar(1024) DEFAULT NULL," +
			"`scheme` varchar(16) DEFAULT NULL," +
			"`proto` varchar(16) DEFAULT NULL," +
			"`bytesSent` bigint(20) unsigned DEFAULT '0'," +
			"`bodyBytesSent` bigint(20) unsigned DEFAULT '0'," +
			"`status` int(11) unsigned DEFAULT '0'," +
			"`statusMessage` varchar(1024) DEFAULT NULL," +
			"`sentHeader` json DEFAULT NULL," +
			"`timeISO8601` varchar(128) DEFAULT NULL," +
			"`timeLocal` varchar(128) DEFAULT NULL," +
			"`msec` decimal(20,6) unsigned DEFAULT '0.000000'," +
			"`timestamp` int(11) unsigned DEFAULT '0'," +
			"`host` varchar(128) DEFAULT NULL," +
			"`referer` varchar(1024) DEFAULT NULL," +
			"`userAgent` varchar(1024) DEFAULT NULL," +
			"`request` varchar(1024) DEFAULT NULL," +
			"`contentType` varchar(256) DEFAULT NULL," +
			"`cookie` json DEFAULT NULL," +
			"`arg` json DEFAULT NULL," +
			"`args` text," +
			"`queryString` text," +
			"`header` json DEFAULT NULL," +
			"`serverName` varchar(256) DEFAULT NULL," +
			"`serverPort` int(11) unsigned DEFAULT '0'," +
			"`serverProtocol` varchar(16) DEFAULT NULL," +
			"`backendAddress` varchar(256) DEFAULT NULL," +
			"`fastcgiAddress` varchar(256) DEFAULT NULL," +
			"`requestData` blob," +
			"`responseHeaderData` blob," +
			"`responseBodyData` blob," +
			"`errors` json DEFAULT NULL," +
			"`hasErrors` tinyint(1) unsigned DEFAULT '0'," +
			"`extend` json DEFAULT NULL," +
			"`attrs` json DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"UNIQUE KEY `_id` (`_id`)," +
			"KEY `serverId` (`serverId`)," +
			"KEY `serverId_status` (`serverId`,`status`)," +
			"KEY `serverId_remoteAddr` (`serverId`,`remoteAddr`)," +
			"KEY `serverId_hasErrors` (`serverId`,`hasErrors`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
		_, err = conn.ExecContext(context.Background(), s)
		if err != nil {
			logs.Error(err)
		}
	}
}
