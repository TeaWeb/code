package teadb

import (
	"context"
	"database/sql"
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/iwind/TeaGo/logs"
	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	SQLDriver
}

// 初始化
func (this *PostgresDriver) Init() {
	this.driver = "postgres"

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

// 初始化数据库
func (this *PostgresDriver) initDB() error {
	config, err := db.LoadPostgresConfig()
	if err != nil {
		return err
	}
	dbInstance, err := sql.Open("postgres", config.DSN)
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
func (this *PostgresDriver) CheckTableExists(table string) (bool, error) {
	conn, err := this.connect()
	if err != nil {
		return false, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	stmt, err := conn.PrepareContext(context.Background(), "SELECT table_name FROM information_schema.tables  WHERE table_name=$1")
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(table)
	i := interface{}(nil)
	err = row.Scan(&i)
	return err == nil, nil
}

// 创建表
func (this *PostgresDriver) CreateTable(table string, definitionSQL string) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}

	exists, err := this.CheckTableExists(table)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	_, err = conn.ExecContext(context.Background(), definitionSQL)
	return err
}
