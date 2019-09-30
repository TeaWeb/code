package teadb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	_ "github.com/lib/pq"
	"net/url"
	"strings"
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
	currentDB, err := this.checkDB()
	if err != nil {
		return false, err
	}

	stmt, err := currentDB.PrepareContext(context.Background(), "SELECT table_name FROM information_schema.tables  WHERE table_name=$1")
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
	currentDB, err := this.checkDB()
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

	this.dbLocker.Lock()
	defer this.dbLocker.Unlock()

	_, err = currentDB.ExecContext(context.Background(), definitionSQL)
	return err
}

// 测试DSN
func (this *PostgresDriver) TestDSN(dsn string, autoCreateDB bool) (message string, ok bool) {
	dbInstance, err := sql.Open("postgres", dsn)
	if err != nil {
		message = "DSN解析错误：" + err.Error()
		return
	}
	defer func() {
		_ = dbInstance.Close()
	}()

	if autoCreateDB {
		u, err := url.Parse(dsn)
		if err != nil {
			message = err.Error()
			return
		}
		if len(u.Path) <= 1 {
			message = "database name should not be empty"
			return
		}
		database := u.Path[1:]
		u.Path = "/"

		newDBInstance, err := sql.Open("postgres", u.String())
		if err != nil {
			message = err.Error()
			return
		}
		_, err = newDBInstance.ExecContext(context.Background(), `CREATE DATABASE "`+database+`"`)
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				message = err.Error()
				_ = newDBInstance.Close()
				return
			}
		}

		_ = newDBInstance.Close()
	}

	// 测试创建数据表
	_, err = dbInstance.ExecContext(context.Background(), `CREATE TABLE "public"."teaweb_test" (
		"id" serial8 primary key,
		"_id" varchar(24)
		);
		
		CREATE UNIQUE INDEX "teaweb_test_id" ON "public"."teaweb_test" ("_id");`)
	if err != nil {
		message = "尝试创建数据表失败：" + err.Error()
		return
	}

	// 测试写入数据表
	_, err = dbInstance.ExecContext(context.Background(), "INSERT INTO \"teaweb_test\" (\"_id\") VALUES ('"+shared.NewObjectId().Hex()+"')")
	if err != nil {
		message = "尝试写入数据表失败：" + err.Error()
		return
	}

	// 测试删除数据表
	_, err = dbInstance.ExecContext(context.Background(), "DROP TABLE \"teaweb_test\"")
	if err != nil {
		message = "尝试删除数据表失败：" + err.Error()
		return
	}

	// 检查函数
	rows, err := dbInstance.Query(`SELECT JSON_EXTRACT_PATH('{"a":1}', 'a')`)
	if err != nil {
		message = "检查JSON_EXTRACT_PATH()函数失败：" + err.Error() + "。请尝试使用PostgreSQL v9.3以上版本。"
		return
	}
	_ = rows.Close()

	ok = true
	return
}

// 统计数据表
// TODO
func (this *PostgresDriver) StatTables(tables []string) (map[string]*TableStat, error) {
	return map[string]*TableStat{}, errors.New("not implemented yet")
}
