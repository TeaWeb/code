package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/db"
	"sync"
)

var (
	sharedDriver   Driver                  = nil
	accessLogDAO   AccessLogDAOInterface   = nil
	agentLogDAO    AgentLogDAOInterface    = nil
	auditLogDAO    AuditLogDAOInterface    = nil
	noticeDAO      NoticeDAOInterface      = nil
	agentValueDAO  AgentValueDAOInterface  = nil
	serverValueDAO ServerValueDAOInterface = nil

	initTableMap    = map[string]bool{}
	initTableLocker = sync.Mutex{}
)

func SetupDB() {
	dbConfig := db.SharedDBConfig()
	switch dbConfig.Type {
	case db.DBTypeMongo:
		sharedDriver = new(MongoDriver)
		/**case db.DBTypeMySQL:
			sharedDriver = new(MySQLDriver)
		case db.DBTypePostgres:
			sharedDriver = new(PostgresDriver)**/
	}

	// initialize
	if sharedDriver != nil {
		sharedDriver.Init()
	}
}

func SharedDB() Driver {
	return sharedDriver
}

func AccessLogDAO() AccessLogDAOInterface {
	return accessLogDAO
}

func AgentLogDAO() AgentLogDAOInterface {
	return agentLogDAO
}

func AuditLogDAO() AuditLogDAOInterface {
	return auditLogDAO
}

func NoticeDAO() NoticeDAOInterface {
	return noticeDAO
}

func AgentValueDAO() AgentValueDAOInterface {
	return agentValueDAO
}

func ServerValueDAO() ServerValueDAOInterface {
	return serverValueDAO
}

func isInitializedTable(table string) bool {
	initTableLocker.Lock()
	defer initTableLocker.Unlock()

	_, ok := initTableMap[table]
	if ok {
		return true
	}

	initTableMap[table] = true
	return false
}
