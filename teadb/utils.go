package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"sync"
	"time"
)

var (
	sharedDriver   DriverInterface         = nil
	accessLogDAO   AccessLogDAOInterface   = nil
	agentLogDAO    AgentLogDAOInterface    = nil
	auditLogDAO    AuditLogDAOInterface    = nil
	noticeDAO      NoticeDAOInterface      = nil
	agentValueDAO  AgentValueDAOInterface  = nil
	serverValueDAO ServerValueDAOInterface = nil

	initTableMap    = map[string]bool{}
	initTableLocker = sync.Mutex{}
)

// 建立数据库驱动
func SetupDB() {
	dbConfig := db.SharedDBConfig()
	switch dbConfig.Type {
	case db.DBTypeMongo:
		sharedDriver = new(MongoDriver)
	case db.DBTypeMySQL:
		sharedDriver = new(MySQLDriver)
		/**case db.DBTypePostgres:
		sharedDriver = new(PostgresDriver)**/
	}

	// initialize
	if sharedDriver != nil {
		sharedDriver.Init()
		sharedDriver.SetIsAvailable(true)

		// 测试数据库
		timers.Loop(10*time.Second, func(looper *timers.Looper) {
			err := sharedDriver.Test()
			if err != nil {
				logs.Println("[db]database connection unavailable: " + err.Error())
			}
			sharedDriver.SetIsAvailable(err == nil)
		})
	}
}

func SharedDB() DriverInterface {
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
