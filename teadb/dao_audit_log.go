package teadb

import "github.com/TeaWeb/code/teaconfigs/audits"

type AuditLogDAOInterface interface {
	Init()
	CountAllAuditLogs() (int64, error)
	ListAuditLogs(offset int, size int) ([]*audits.Log, error)
	InsertOne(auditLog *audits.Log) error
}
