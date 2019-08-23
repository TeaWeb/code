package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/db"
)

var sharedDriver Driver = nil

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
