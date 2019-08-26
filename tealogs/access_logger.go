package tealogs

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/code/tealogs/accesslogs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"runtime"
)

var (
	accessLogger *AccessLogger = nil
)

// 访问日志记录器
type AccessLogger struct {
	queue chan *accesslogs.AccessLog
}

// 获取新日志对象
func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue: make(chan *accesslogs.AccessLog, 10*10000),
	}

	go logger.wait()
	return logger
}

// 获取共享的对象
func SharedLogger() *AccessLogger {
	return accessLogger
}

// 推送日志
func (this *AccessLogger) Push(log *accesslogs.AccessLog) {
	this.queue <- log
}

// 等待日志到来
func (this *AccessLogger) wait() {
	var logDBNames = []string{}
	for i := 0; i < runtime.NumCPU(); i++ {
		if i == 0 {
			logDBNames = append(logDBNames, "accesslog.leveldb")
		} else {
			logDBNames = append(logDBNames, "accesslog"+fmt.Sprintf("%d", i)+".leveldb")
		}
	}

	// 清除日志
	for _, dbName := range logDBNames {
		logFile := files.NewFile(Tea.Root + "/logs/" + dbName)
		if logFile.Exists() && logFile.IsDir() {
			err := logFile.DeleteAll()
			if err != nil {
				logs.Error(err)
			}
		}
	}

	// 打开leveldb数据库
	var logDBs = []*leveldb.DB{}
	for _, dbName := range logDBNames {
		db, err := leveldb.OpenFile(Tea.Root+"/logs/"+dbName, nil)
		if err != nil {
			logs.Error(errors.New("open " + dbName + ": " + err.Error()))
			return
		}
		logDBs = append(logDBs, db)
	}

	// 启动queue
	for index, db := range logDBs {
		func(index int, db *leveldb.DB) {
			queue := NewAccessLogQueue(db, index)
			go queue.Receive(this.queue)
			go queue.Dump()
		}(index, db)
	}
}
