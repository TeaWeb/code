package tealogs

import (
	"context"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"runtime"
	"sync"
	"time"
)

var (
	accessLogger = NewAccessLogger()
)

// 访问日志记录器
type AccessLogger struct {
	queue chan *AccessLogItem

	logs            []*AccessLogItem
	timestamp       int64
	qps             int
	outputBandWidth int64
	inputBandWidth  int64

	collectionCacheMap    map[string]*mongo.Collection
	collectionCacheLocker sync.Mutex
}

type AccessLogItem struct {
	log *AccessLog
}

func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue:              make(chan *AccessLogItem, 10240),
		collectionCacheMap: map[string]*mongo.Collection{},
	}

	go logger.wait()
	return logger
}

func SharedLogger() *AccessLogger {
	return accessLogger
}

func (this *AccessLogger) Push(log *AccessLog) {
	this.queue <- &AccessLogItem{
		log: log,
	}
}

func (this *AccessLogger) client() *mongo.Client {
	return teamongo.SharedClient()
}

func (this *AccessLogger) collection() *mongo.Collection {
	collName := "logs." + timeutil.Format("Ymd")
	coll, found := this.collectionCacheMap[collName]
	if found {
		return coll
	}

	// 构建索引
	coll = this.client().Database("teaweb").Collection(collName)
	indexes := coll.Indexes()
	indexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.NewDocument(bson.EC.Int32("serverId", 1)),
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})
	indexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.NewDocument(bson.EC.Int32("status", 1), bson.EC.Int32("serverId", 1)),
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})
	indexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.NewDocument(bson.EC.Int32("remoteAddr", 1), bson.EC.Int32("serverId", 1)),
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})
	indexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.NewDocument(bson.EC.Int32("apiPath", 1), bson.EC.Int32("serverId", 1)),
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})

	this.collectionCacheLocker.Lock()
	this.collectionCacheMap[collName] = coll
	this.collectionCacheLocker.Unlock()

	return coll
}

func (this *AccessLogger) wait() {
	timestamp := time.Now().Unix()

	var docs = []interface{}{}
	var docsLocker = sync.Mutex{}

	// 写入到数据库
	countCPU := runtime.NumCPU()
	timers.Loop(500*time.Millisecond, func(looper *timers.Looper) {
		// 写入到本地数据库
		if this.client() != nil {
			docsLocker.Lock()
			if len(docs) == 0 {
				docsLocker.Unlock()
				return
			}
			newDocs := docs
			docs = []interface{}{}
			docsLocker.Unlock()

			total := len(newDocs)

			// 批量写入数据库
			// 需合理控制此数值的大小，避免CPU占用太高
			bulkSize := countCPU * 64
			offset := 0
			for {
				end := offset + bulkSize
				if end > total {
					end = total
				}

				logs.Println("dump", end-offset, "access logs ...")
				docSlice := newDocs[offset:end]

				// 分析
				for _, doc := range docSlice {
					accessLog := doc.(*AccessLog)
					accessLog.Parse()
					accessLog.Id = objectid.New()

					// 执行处理器
					CallAccessLogHooks(accessLog)
				}

				// 写入数据库
				_, err := this.collection().InsertMany(context.Background(), docSlice)
				if err != nil {
					logs.Error(err)
					return
				}
				logs.Println("done")
				docSlice = []interface{}{}

				offset = end
				if end >= total {
					break
				}
			}

			// 清空
			newDocs = []interface{}{}
		}
	})

	// 接收日志
	for {
		item := <-this.queue
		log := item.log

		// 计算QPS和BandWidth
		this.timestamp = log.Timestamp
		if log.Timestamp == timestamp {
			this.qps ++
			this.inputBandWidth += log.RequestLength
			this.outputBandWidth += log.BytesSent
		} else {
			this.qps = 1
			this.inputBandWidth = log.RequestLength
			this.outputBandWidth = log.BytesSent
			timestamp = log.Timestamp
		}

		docsLocker.Lock()
		docs = append(docs, log)
		docsLocker.Unlock()
	}
}

// 关闭
func (this *AccessLogger) Close() {
	if this.client() != nil {
		this.client().Disconnect(context.Background())
	}
}

// 读取日志
func (this *AccessLogger) ReadNewLogs(serverId string, fromId string, size int64) []AccessLog {
	if this.client() == nil {
		return []AccessLog{}
	}

	if size <= 0 {
		size = 10
	}

	result := []AccessLog{}
	coll := this.collection()

	filter := map[string]interface{}{}
	if len(serverId) > 0 {
		filter["serverId"] = serverId
	}
	if len(fromId) > 0 {
		objectId, err := objectid.FromHex(fromId)
		if err == nil {
			filter["_id"] = map[string]interface{}{
				"$gt": objectId,
			}
		} else {
			logs.Error(err)
		}
	}

	opts := []findopt.Find{}
	isReverse := false

	if len(fromId) == 0 {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("_id", -1))))
		opts = append(opts, findopt.Limit(size))
		isReverse = true
	} else {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("_id", 1))))
		opts = append(opts, findopt.Limit(size))
	}

	cursor, err := coll.Find(context.Background(), filter, opts ...)
	if err != nil {
		logs.Error(err)
		return []AccessLog{}
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		accessLog := AccessLog{}
		err := cursor.Decode(&accessLog)
		if err != nil {
			logs.Error(err)
			return []AccessLog{}
		}
		result = append(result, accessLog)
	}

	if !isReverse {
		lists.Reverse(result)
	}
	return result
}

func (this *AccessLogger) QPS() int {
	if time.Now().Unix()-this.timestamp < 2 {
		return this.qps
	}
	return 0
}

func (this *AccessLogger) InputBandWidth() int64 {
	if time.Now().Unix()-this.timestamp < 2 {
		return this.inputBandWidth
	}
	return 0
}

func (this *AccessLogger) OutputBandWidth() int64 {
	if time.Now().Unix()-this.timestamp < 2 {
		return this.outputBandWidth
	}
	return 0
}

// 读取日志
func (this *AccessLogger) ReadNewLogsForAPI(serverId string, apiPath string, fromId string, size int64) []AccessLog {
	if this.client() == nil {
		return []AccessLog{}
	}

	if size <= 0 {
		size = 10
	}

	result := []AccessLog{}
	coll := this.collection()

	filter := map[string]interface{}{
		"serverId": serverId,
		"apiPath":  apiPath,
	}
	if len(fromId) > 0 {
		objectId, err := objectid.FromHex(fromId)
		if err == nil {
			filter["_id"] = map[string]interface{}{
				"$gt": objectId,
			}
		} else {
			logs.Error(err)
		}
	}

	opts := []findopt.Find{}
	isReverse := false

	if len(fromId) == 0 {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("_id", -1))))
		opts = append(opts, findopt.Limit(size))
		isReverse = true
	} else {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("_id", 1))))
		opts = append(opts, findopt.Limit(size))
	}

	cursor, err := coll.Find(context.Background(), filter, opts ...)
	if err != nil {
		logs.Error(err)
		return []AccessLog{}
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		accessLog := AccessLog{}
		err := cursor.Decode(&accessLog)
		if err != nil {
			logs.Error(err)
			return []AccessLog{}
		}
		result = append(result, accessLog)
	}

	if !isReverse {
		lists.Reverse(result)
	}
	return result
}
