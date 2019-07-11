package tealogs

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/mongo"
	"runtime"
	"sync"
)

var (
	accessLogger *AccessLogger = nil
)

// 访问日志记录器
type AccessLogger struct {
	queue chan *AccessLog

	collectionCacheMap    map[string]*teamongo.Collection
	collectionCacheLocker sync.RWMutex
}

// 获取新日志对象
func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue:              make(chan *AccessLog, 10*10000),
		collectionCacheMap: map[string]*teamongo.Collection{},
	}

	go logger.wait()
	return logger
}

// 获取共享的对象
func SharedLogger() *AccessLogger {
	return accessLogger
}

// 推送日志
func (this *AccessLogger) Push(log *AccessLog) {
	this.queue <- log
}

// 获取MongoDB客户端
func (this *AccessLogger) client() *mongo.Client {
	return teamongo.SharedClient()
}

// 获取当天的collection
func (this *AccessLogger) collection() *teamongo.Collection {
	this.collectionCacheLocker.RLock()

	collName := "logs." + timeutil.Format("Ymd")
	coll, found := this.collectionCacheMap[collName]
	if found {
		this.collectionCacheLocker.RUnlock()
		return coll
	}
	this.collectionCacheLocker.RUnlock()

	// 构建索引
	this.collectionCacheLocker.Lock()
	coll = teamongo.FindCollection(collName)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
	})
	coll.CreateIndex(map[string]bool{
		"status":   true,
		"serverId": true,
	})
	coll.CreateIndex(map[string]bool{
		"remoteAddr": true,
		"serverId":   true,
	})
	coll.CreateIndex(map[string]bool{
		"hasErrors": true,
		"serverId":  true,
	})

	this.collectionCacheMap[collName] = coll
	this.collectionCacheLocker.Unlock()
	return coll
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
			go queue.Dump(this.collection)
		}(index, db)
	}
}

// 关闭MongoDB客户端连接
func (this *AccessLogger) Close() {
	if this.client() != nil {
		this.client().Disconnect(context.Background())
	}
}
