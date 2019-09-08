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

// 获取共享的数据库驱动
func SharedDB() DriverInterface {
	return sharedDriver
}

// 获取访问日志DAO
func AccessLogDAO() AccessLogDAOInterface {
	return accessLogDAO
}

// 获取Agent日志DAO
func AgentLogDAO() AgentLogDAOInterface {
	return agentLogDAO
}

// 获取审计日志DAO
func AuditLogDAO() AuditLogDAOInterface {
	return auditLogDAO
}

// 获取通知DAO
func NoticeDAO() NoticeDAOInterface {
	return noticeDAO
}

// 获取Agent数值记录DAO
func AgentValueDAO() AgentValueDAOInterface {
	return agentValueDAO
}

// 获取代理统计数值DAO
func ServerValueDAO() ServerValueDAOInterface {
	return serverValueDAO
}

// 判断表格是否已经初始化
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
