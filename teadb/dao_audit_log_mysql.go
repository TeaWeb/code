package teadb

import (
	"context"
	"github.com/TeaWeb/code/teaconfigs/audits"
	"github.com/iwind/TeaGo/logs"
)

type MySQLAuditLogDAO struct {
}

// 初始化
func (this *MySQLAuditLogDAO) Init() {

}

func (this *MySQLAuditLogDAO) TableName() string {
	this.initTable("teaweb.logs.audit")
	return "teaweb.logs.audit"
}

// 计算审计日志数量
func (this *MySQLAuditLogDAO) CountAllAuditLogs() (int64, error) {
	return NewQuery(this.TableName()).Count()
}

// 列出审计日志
func (this *MySQLAuditLogDAO) ListAuditLogs(offset int, size int) ([]*audits.Log, error) {
	ones, err := NewQuery(this.TableName()).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(audits.Log))
	if err != nil {
		return nil, err
	}
	result := []*audits.Log{}
	for _, one := range ones {
		result = append(result, one.(*audits.Log))
	}
	return result, nil
}

// 插入一条审计日志
func (this *MySQLAuditLogDAO) InsertOne(auditLog *audits.Log) error {
	return SharedDB().InsertOne(this.TableName(), auditLog)
}

// 初始化表格
func (this *MySQLAuditLogDAO) initTable(table string) {
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
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID'," +
			"`_id` varchar(24) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'global id'," +
			"`action` varchar(256) COLLATE utf8mb4_bin DEFAULT NULL," +
			"`username` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL," +
			"`description` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL," +
			"`options` json DEFAULT NULL," +
			"`timestamp` int(11) DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"UNIQUE KEY `_id` (`_id`)" +
			") ENGINE=InnoDB AUTO_INCREMENT=100013 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;"
		_, err = conn.ExecContext(context.Background(), s)
		if err != nil {
			logs.Error(err)
		}
	}
}
