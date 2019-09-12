package teadb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/TeaWeb/code/teadb/shared"
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
func (this *MySQLDriver) TestDSN(dsn string) (message string, ok bool) {
	dbInstance, err := sql.Open("mysql", dsn)
	if err != nil {
		message = "DSN解析错误：" + err.Error()
		return
	}

	conn, err := dbInstance.Conn(context.Background())
	if err != nil {
		message = "尝试连接数据库失败：" + err.Error()
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	// 测试创建数据表
	_, err = conn.ExecContext(context.Background(), "CREATE TABLE `teaweb.test` ( "+
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
	_, err = conn.ExecContext(context.Background(), "INSERT INTO `teaweb.test` (`_id`) VALUES (\""+shared.NewObjectId().Hex()+"\")")
	if err != nil {
		message = "尝试写入数据表失败：" + err.Error()
		return
	}

	// 测试删除数据表
	_, err = conn.ExecContext(context.Background(), "DROP TABLE `teaweb.test`")
	if err != nil {
		message = "尝试删除数据表失败：" + err.Error()
		return
	}

	ok = true
	return
}

// 统计数据表
// TODO
func (this *MySQLDriver) StatTables(tables []string) (map[string]*TableStat, error) {
	return map[string]*TableStat{}, errors.New("not implemented yet")
}
