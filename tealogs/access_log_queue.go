package tealogs

import (
	"context"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"sync"
	"time"
)

// 访问日志队列
type AccessLogQueue struct {
	index  int
	db     *leveldb.DB
	locker sync.Mutex
	logIds []string
}

// 创建队列对象
func NewAccessLogQueue(db *leveldb.DB, index int) *AccessLogQueue {
	return &AccessLogQueue{
		db:    db,
		index: index,
	}
}

// 从队列中接收日志
func (this *AccessLogQueue) Receive(ch chan *AccessLog) {
	ticker := time.NewTicker(1 * time.Second)
	batch := new(leveldb.Batch)
	logIds := []string{}
	id := uint64(time.Now().UnixNano())

	for {
		select {
		case log := <-ch:
			if log != nil {
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
						idString := strconv.FormatUint(id, 10)
						logIds = append(logIds, idString)
						batch.Put([]byte("accesslog_"+idString), data)
					}
				}

				// 批量写入
				if batch.Len() > 2048 {
					err := this.db.Write(batch, nil)
					if err != nil {
						logs.Error(err)
					}
					batch.Reset()
					this.locker.Lock()
					this.logIds = append(this.logIds, logIds...)
					this.locker.Unlock()
					logIds = []string{}
				}
			}
		case <-ticker.C:
			if batch.Len() > 0 {
				err := this.db.Write(batch, nil)
				if err != nil {
					logs.Error(err)
				}
				batch.Reset()
				this.locker.Lock()
				this.logIds = append(this.logIds, logIds...)
				this.locker.Unlock()
				logIds = []string{}
			}
		}
	}
}

// 导出日志到别的媒介
func (this *AccessLogQueue) Dump(mongoCollFunc func() *teamongo.Collection) {
	for {
		this.dumpInterval(mongoCollFunc)
		time.Sleep(1 * time.Second)
	}
}

// 导出日志定时内容
func (this *AccessLogQueue) dumpInterval(mongoCollFunc func() *teamongo.Collection) {
	size := 4096
	var logIds = []string{}
	this.locker.Lock()

	// 超出一定数值，则清空，防止占用内存过大
	if len(this.logIds) > 100*10000 {
		this.logIds = []string{}
	} else if len(this.logIds) > size {
		logIds = this.logIds[:size]
		this.logIds = this.logIds[size:]
	} else {
		logIds = this.logIds
		this.logIds = []string{}
	}
	this.locker.Unlock()

	accessLogs := []interface{}{}
	batch := new(leveldb.Batch)

	for _, logId := range logIds {
		key := []byte("accesslog_" + logId)
		value, err := this.db.Get(key, nil)
		if err != nil {
			if err != leveldb.ErrNotFound {
				logs.Error(err)
				this.db.Delete(key, nil)
			}
			continue
		}
		accessLog := new(AccessLog)
		err = ffjson.Unmarshal(value, accessLog)
		if err != nil {
			logs.Error(err)
			this.db.Delete(key, nil)
			continue
		}
		accessLog.Id = primitive.NewObjectID()
		accessLogs = append(accessLogs, accessLog)

		batch.Delete(key)
	}

	if batch.Len() > 0 {
		this.db.Write(batch, nil)
	}

	if len(accessLogs) > 0 {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := mongoCollFunc().InsertMany(ctx, accessLogs)
		if err != nil {
			logs.Println("[mongo]insert access logs:", err.Error())
		}
	}
}
