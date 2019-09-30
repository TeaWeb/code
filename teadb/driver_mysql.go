package teadb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/TeaWeb/code/teadb/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/logs"
	"strings"
)

type MySQLDriver struct {
	SQLDriver
}

// 初始化
func (this *MySQLDriver) Init() {
	this.driver = "mysql"

	err := this.initDB()
	if err != nil {
		logs.Error(err)
	}

	agentValueDAO = new(SQLAgentValueDAO)
	agentValueDAO.SetDriver(this)
	agentValueDAO.Init()

	agentLogDAO = new(SQLAgentLogDAO)
	agentLogDAO.SetDriver(this)
	agentLogDAO.Init()

	serverValueDAO = new(SQLServerValueDAO)
	serverValueDAO.SetDriver(this)
	serverValueDAO.Init()

	auditLogDAO = new(SQLAuditLogDAO)
	auditLogDAO.SetDriver(this)
	auditLogDAO.Init()

	accessLogDAO = new(SQLAccessLogDAO)
	accessLogDAO.SetDriver(this)
	accessLogDAO.Init()

	noticeDAO = new(SQLNoticeDAO)
	noticeDAO.SetDriver(this)
	noticeDAO.Init()
}

func (this *MySQLDriver) initDB() error {
	config, err := db.LoadMySQLConfig()
	if err != nil {
		return err
	}
	dbInstance, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return err
	}
	dbInstance.SetMaxIdleConns(10)
	dbInstance.SetMaxOpenConns(32)
	dbInstance.SetConnMaxLifetime(0)
	this.db = dbInstance

	go func() {
		row := this.db.QueryRow("SELECT @@sql_mode")
		if row != nil {
			_ = row.Scan(&this.sqlMode)
		}
	}()

	return nil
}

// 检查表是否存在
func (this *MySQLDriver) CheckTableExists(table string) (bool, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return false, err
	}

	_, err = currentDB.ExecContext(context.Background(), "SHOW CREATE TABLE `"+table+"`")
	if err != nil {
		return false, nil
	}

	return true, nil
}

// 创建表
func (this *MySQLDriver) CreateTable(table string, definitionSQL string) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}

	exists, err := this.CheckTableExists(table)
	if err != nil {
		return err
	}
	if !exists {
		_, err = currentDB.ExecContext(context.Background(), definitionSQL)
		if err != nil {
			logs.Error(err)
		}
	}
	return err
}

// 测试DSN
func (this *MySQLDriver) TestDSN(dsn string, autoCreateDB bool) (message string, ok bool) {
	dbInstance, err := sql.Open("mysql", dsn)
	if err != nil {
		message = "DSN解析错误：" + err.Error()
		return
	}
	defer func() {
		_ = dbInstance.Close()
	}()

	// 检查数据库
	if autoCreateDB {
		index := strings.Index(dsn, "/")
		if index == -1 {
			message = "invalid dsn"
			return
		}
		database := dsn[index+1:]
		index = strings.Index(database, "?")
		if index > -1 {
			database = database[:index]
		}
		if len(database) == 0 {
			message = "no database defined"
			return
		}
		newDSN := strings.Replace(dsn, "/"+database, "/", -1)
		newDBInstance, err := sql.Open("mysql", newDSN)
		if err != nil {
			message = err.Error()
			return
		}
		_, err = newDBInstance.ExecContext(context.Background(), "CREATE DATABASE IF NOT EXISTS `"+database+"`")
		if err != nil {
			message = err.Error()
			_ = newDBInstance.Close()
			return
		}
		_ = newDBInstance.Close()
	}

	// 测试创建数据表
	_, err = dbInstance.ExecContext(context.Background(), "CREATE TABLE `teaweb_test` ( "+
		"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,"+
		"`_id` varchar(24) DEFAULT NULL,"+
		"PRIMARY KEY (`id`),"+
		"UNIQUE KEY `_id` (`_id`)"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		message = "尝试创建数据表失败：" + err.Error()
		return
	}

	// 测试写入数据表
	_, err = dbInstance.ExecContext(context.Background(), "INSERT INTO `teaweb_test` (`_id`) VALUES (\""+shared.NewObjectId().Hex()+"\")")
	if err != nil {
		message = "尝试写入数据表失败：" + err.Error()
		return
	}

	// 测试删除数据表
	_, err = dbInstance.ExecContext(context.Background(), "DROP TABLE `teaweb_test`")
	if err != nil {
		message = "尝试删除数据表失败：" + err.Error()
		return
	}

	// 检查函数
	rows, err := dbInstance.Query(`SELECT JSON_EXTRACT('{"a":1}', "$.a")`)
	if err != nil {
		message = "检查JSON_EXTRACT()函数失败：" + err.Error() + "。请尝试使用MySQL v5.7.8以上版本。"
		return
	}
	_ = rows.Close()

	ok = true
	return
}

// 统计数据表
// TODO
func (this *MySQLDriver) StatTables(tables []string) (map[string]*TableStat, error) {
	return map[string]*TableStat{}, errors.New("not implemented yet")
}
