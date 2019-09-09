package teadb

import (
	"context"
	"database/sql"
	"github.com/TeaWeb/code/teaconfigs/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/logs"
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
	agentValueDAO.Init()

	agentLogDAO = new(SQLAgentLogDAO)
	agentLogDAO.Init()

	serverValueDAO = new(SQLServerValueDAO)
	serverValueDAO.Init()

	auditLogDAO = new(SQLAuditLogDAO)
	auditLogDAO.Init()

	accessLogDAO = new(SQLAccessLogDAO)
	accessLogDAO.Init()

	noticeDAO = new(SQLNoticeDAO)
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
	return nil
}

// 检查表是否存在
func (this *MySQLDriver) CheckTableExists(table string) (bool, error) {
	conn, err := this.connect()
	if err != nil {
		return false, err
	}

	_, err = conn.ExecContext(context.Background(), "SHOW CREATE TABLE `"+table+"`")
	if err != nil {
		return false, nil
	}

	return true, nil
}

// 创建表
func (this *MySQLDriver) CreateTable(table string, definitionSQL string) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}

	exists, err := this.CheckTableExists(table)
	if err != nil {
		return err
	}
	if !exists {
		_, err = conn.ExecContext(context.Background(), definitionSQL)
		if err != nil {
			logs.Error(err)
		}
	}
	return err
}
