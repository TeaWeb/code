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
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		this.dumpInterval(mongoCollFunc)
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

	if len(logIds) == 0 {
		return
	}

	accessLogs := []interface{}{}
	batch := new(leveldb.Batch)

	storageLogs := map[string][]*AccessLog{} // policyId => accessLogs
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

		// 如果非storageOnly则可以存储到MongoDB中
		if !accessLog.StorageOnly {
			accessLog.Id = primitive.NewObjectID()
			accessLogs = append(accessLogs, accessLog)
		}

		// 日志存储策略
		if len(accessLog.StoragePolicyIds) > 0 {
			for _, policyId := range accessLog.StoragePolicyIds {
				_, ok := storageLogs[policyId]
				if !ok {
					storageLogs[policyId] = []*AccessLog{}
				}
				storageLogs[policyId] = append(storageLogs[policyId], accessLog)
			}
		}

		batch.Delete(key)
	}

	if len(storageLogs) > 0 {
		for policyId, storageAccessLogs := range storageLogs {
			storage := FindPolicyStorage(policyId)
			if storage == nil {
				continue
			}
			err := storage.Write(storageAccessLogs)
			if err != nil {
				logs.Println("access log storage policy '"+policyId+"/"+FindPolicyName(policyId)+"'", err.Error())
			}
		}
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
