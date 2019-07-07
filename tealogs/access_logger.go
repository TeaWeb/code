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
	collectionCacheLocker sync.Mutex
}

func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue:              make(chan *AccessLog, 10*10000),
		collectionCacheMap: map[string]*teamongo.Collection{},
	}

	go logger.wait()
	return logger
}

func SharedLogger() *AccessLogger {
	return accessLogger
}

func (this *AccessLogger) Push(log *AccessLog) {
	this.queue <- log
}

func (this *AccessLogger) client() *mongo.Client {
	return teamongo.SharedClient()
}

func (this *AccessLogger) collection() *teamongo.Collection {
	this.collectionCacheLocker.Lock()
	defer this.collectionCacheLocker.Unlock()

	collName := "logs." + timeutil.Format("Ymd")
	coll, found := this.collectionCacheMap[collName]
	if found {
		return coll
	}

	// 构建索引
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

	var logDBs = []*leveldb.DB{}
	for _, dbName := range logDBNames {
		db, err := leveldb.OpenFile(Tea.Root+"/logs/"+dbName, nil)
		if err != nil {
			logs.Error(errors.New("open " + dbName + ": " + err.Error()))
			return
		}
		logDBs = append(logDBs, db)
	}

	for index, db := range logDBs {
		func(index int, db *leveldb.DB) {
			queue := NewAccessLogQueue(db, index)
			go queue.Receive(this.queue)
			go queue.Dump(this.collection)
		}(index, db)
	}
}

// 关闭
func (this *AccessLogger) Close() {
	if this.client() != nil {
		this.client().Disconnect(context.Background())
	}
}
