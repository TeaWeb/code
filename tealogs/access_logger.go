package tealogs

import (
	"context"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	accessLogger *AccessLogger = nil
	globalLogIds               = []string{}
	logLocker                  = &sync.Mutex{}
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
	// 接收日志
	logFile := files.NewFile(Tea.Root + "/logs/accesslog.leveldb")
	if logFile.Exists() && logFile.IsDir() {
		err := logFile.DeleteAll()
		if err != nil {
			logs.Error(err)
		}
	}
	logsDb, err := leveldb.OpenFile(Tea.Root+"/logs/accesslog.leveldb", nil)
	if err != nil {
		logs.Error(err)
		return
	}

	id := uint64(time.Now().UnixNano())

	for i := 0; i < runtime.NumCPU(); i++ {
		go this.receive(logsDb, &id)
	}

	for {
		this.dumpLogsToMongo(logsDb)
		time.Sleep(1 * time.Second)
	}
}

// 关闭
func (this *AccessLogger) Close() {
	if this.client() != nil {
		this.client().Disconnect(context.Background())
	}
}

// 接收日志
func (this *AccessLogger) receive(db *leveldb.DB, idPtr *uint64) {
	ticker := time.NewTicker(1 * time.Second)
	batch := new(leveldb.Batch)
	logIds := []string{}
	for log := range this.queue {
		if log == nil {
			break
		}
		newId := atomic.AddUint64(idPtr, 1)

		if log.ShouldStat() || log.ShouldWrite() {
			log.Parse()

			// 统计
			if log.ShouldStat() {
				CallAccessLogHooks(log)
			}

			// 保存到文件
			if log.ShouldWrite() {
				log.CleanFields()
				data, err := ffjson.Marshal(log)
				if err != nil {
					logs.Error(err)
					continue
				}
				idString := strconv.FormatUint(newId, 10)
				logIds = append(logIds, idString)
				batch.Put([]byte("accesslog_"+idString), data)
				select {
				case <-ticker.C:
					if batch.Len() > 0 {
						err = db.Write(batch, nil)
						if err != nil {
							logs.Error(err)
						}
						batch.Reset()
						logLocker.Lock()
						globalLogIds = append(globalLogIds, logIds...)
						logLocker.Unlock()
						logIds = []string{}
					}
				default:
					if batch.Len() > 2048 {
						err = db.Write(batch, nil)
						if err != nil {
							logs.Error(err)
						}
						batch.Reset()
						logLocker.Lock()
						globalLogIds = append(globalLogIds, logIds...)
						logLocker.Unlock()
						logIds = []string{}
					}
				}
			}
		}
	}
}

// 把日志从文件发送到MongoDB
func (this *AccessLogger) dumpLogsToMongo(db *leveldb.DB) {
	size := 4096
	var logIds = []string{}
	logLocker.Lock()

	// 超出一定数值，则清空，防止占用内存过大
	if len(globalLogIds) > 100*10000 {
		globalLogIds = []string{}
	} else if len(globalLogIds) > size {
		logIds = globalLogIds[:size]
		globalLogIds = globalLogIds[size:]
	} else {
		logIds = globalLogIds
		globalLogIds = []string{}
	}
	logLocker.Unlock()

	accessLogs := []interface{}{}
	batch := new(leveldb.Batch)
	for _, logId := range logIds {
		key := []byte("accesslog_" + logId)
		value, err := db.Get(key, nil)
		if err != nil {
			if err != leveldb.ErrNotFound {
				logs.Error(err)
				db.Delete(key, nil)
			}
			continue
		}
		accessLog := new(AccessLog)
		err = ffjson.Unmarshal(value, accessLog)
		if err != nil {
			logs.Error(err)
			db.Delete(key, nil)
			continue
		}
		accessLog.Id = primitive.NewObjectID()
		accessLogs = append(accessLogs, accessLog)

		batch.Delete(key)
	}

	if batch.Len() > 0 {
		db.Write(batch, nil)
	}

	if len(accessLogs) > 0 {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := this.collection().InsertMany(ctx, accessLogs)
		if err != nil {
			logs.Println("[mongo]insert access logs:", err.Error())
		}
	}
}
