package teadb

import "github.com/iwind/TeaGo/maps"

type Driver interface {
	Init()
	FindOne(query *Query, modelPtr interface{}) (interface{}, error)
	FindOnes(query *Query, modelPtr interface{}) ([]interface{}, error)
	InsertOne(table string, modelPtr interface{}) error
	InsertOnes(table string, modelPtrSlice interface{}) error
	DeleteOnes(query *Query) error

	Count(query *Query) (int64, error)
	Sum(query *Query, field string) (float64, error)
	Avg(query *Query, field string) (float64, error)
	Min(query *Query, field string) (float64, error)
	Max(query *Query, field string) (float64, error)

	Group(query *Query, field string, result map[string]Expr) ([]maps.Map, error)

	/**
	UpdateOne(one interface{}) error
	UpdateOnes(ones interface{}) error
	DeleteOne() error
	Drop() error()
	Create() error
	Truncate() error
	**/

	AccessLogDAO() AccessLogDAO
	AgentLogDAO() AgentLogDAO
	AuditLogDAO() AuditLogDAO
	NoticeDAO() NoticeDAO
	ValueDAO() ValueDAO
}
