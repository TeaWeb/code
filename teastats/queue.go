package teastats

import (
	"errors"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"sync"
	"time"
)

// 入库队列
type Queue struct {
	ServerId    string
	coll        *teamongo.Collection
	c           chan *Value
	looper      *timers.Looper
	indexes     [][]string // { { page }, { region, province, city }, ... }
	indexLocker sync.Mutex
}

// 获取新对象
func NewQueue() *Queue {
	return &Queue{}
}

// 队列单例
func (this *Queue) Start(serverId string) {
	this.ServerId = serverId
	this.c = make(chan *Value, 4096)

	// 测试连接，如果有错误则重新连接
	err := teamongo.Test()
	if err != nil {
		logs.Println("[stat]queue start failed: can not connect to mongodb, will reconnect to mongodb")
		time.Sleep(5 * time.Second)
		this.Start(serverId)
		return
	}

	insertQuery := teamongo.NewQuery("values.server."+serverId, new(Value))

	// 创建索引
	coll := insertQuery.Coll()

	this.coll = coll
	for _, indexMap := range []map[string]bool{
		{
			"item":      true,
			"timestamp": true,
		},
		{
			"item":              true,
			"timeFormat.second": true,
		},
		{
			"item":              true,
			"timeFormat.minute": true,
		},
		{
			"item":            true,
			"timeFormat.hour": true,
		},
		{
			"item":           true,
			"timeFormat.day": true,
		},
		{
			"item":            true,
			"timeFormat.week": true,
		},
		{
			"item":             true,
			"timeFormat.month": true,
		},
		{
			"item":            true,
			"timeFormat.year": true,
		},
	} {
		err := coll.CreateIndex(indexMap)
		if err != nil {
			logs.Error(errors.New("mongo:" + err.Error()))
		}
	}

	// 导入数据
	go func() {
		for {
			item := <-this.c

			if item == nil {
				break
			}

			// 是否已存在
			findQuery := teamongo.NewQuery("values.server."+serverId, new(Value)).
				Attr("item", item.Item)
			switch item.Period {
			case ValuePeriodSecond:
				findQuery.Attr("timestamp", item.Timestamp)
			case ValuePeriodMinute:
				findQuery.Attr("timeFormat.minute", item.TimeFormat.Minute)
			case ValuePeriodHour:
				findQuery.Attr("timeFormat.hour", item.TimeFormat.Hour)
			case ValuePeriodDay:
				findQuery.Attr("timeFormat.day", item.TimeFormat.Day)
			case ValuePeriodWeek:
				findQuery.Attr("timeFormat.week", item.TimeFormat.Week)
			case ValuePeriodMonth:
				findQuery.Attr("timeFormat.month", item.TimeFormat.Month)
			case ValuePeriodYear:
				findQuery.Attr("timeFormat.year", item.TimeFormat.Year)
			}

			// 参数
			if len(item.Params) > 0 {
				for k, v := range item.Params {
					findQuery.Attr("params."+k, v)
				}
			} else {
				findQuery.Attr("params", map[string]string{})
			}

			one, err := findQuery.Find()
			if err != nil {
				logs.Error(err)
				continue
			}
			if one == nil {
				// 是否有自定义运算函数
				increaseFunc, found := item.Value["$increase"]
				if found {
					item.Value = increaseFunc.(func(value maps.Map, inc maps.Map) maps.Map)(nil, item.Value)
					delete(item.Value, "$increase")
				}
				err := insertQuery.Insert(item)
				if err != nil {
					logs.Error(err)
				}
			} else {
				oneValue := one.(*Value)

				// 是否有自定义运算函数
				increaseFunc, found := item.Value["$increase"]
				if found {
					item.Value = increaseFunc.(func(value maps.Map, inc maps.Map) maps.Map)(oneValue.Value, item.Value)
					delete(item.Value, "$increase")
				} else {
					// 简单的增长
					item.Value = this.increase(oneValue.Value, item.Value)
				}
				err := teamongo.NewQuery("values.server."+serverId, new(Value)).
					Attr("_id", oneValue.Id).
					Update(map[string]interface{}{
						"$set": maps.Map{
							"value":     item.Value,
							"timestamp": item.Timestamp,
						},
					})
				if err != nil {
					logs.Error(err)
				}
			}
		}
	}()

	// 清理数据
	go func() {
		this.looper = timers.Loop(1*time.Hour, func(looper *timers.Looper) {
			// 清除24小时之前的second
			err := teamongo.NewQuery("values.server."+serverId, new(Value)).
				Attr("period", "second").
				Lt("timestamp", time.Now().Unix()-24*3600).
				Delete()
			if err != nil {
				logs.Error(err)
			}

			// 清除24小时之前的minute
			err = teamongo.NewQuery("values.server."+serverId, new(Value)).
				Attr("period", "minute").
				Lt("timestamp", time.Now().Unix()-24*3600).
				Delete()
			if err != nil {
				logs.Error(err)
			}

			// 清除48小时之前的hour
			err = teamongo.NewQuery("values.server."+serverId, new(Value)).
				Attr("period", "hour").
				Lt("timestamp", time.Now().Unix()-48*3600).
				Delete()
			if err != nil {
				logs.Error(err)
			}
		})
	}()
}

// 添加指标值
func (this *Queue) Add(itemCode string, t time.Time, period ValuePeriod, params map[string]string, value maps.Map) {
	if params == nil {
		params = map[string]string{}
	}
	if value == nil {
		value = maps.Map{}
	}
	item := NewItemValue()
	item.Id = primitive.NewObjectID()
	item.Item = itemCode
	item.Period = period
	item.Value = value
	item.Params = params
	item.SetTime(t)

	if this.c != nil {
		this.c <- item
	}
}

// 停止
func (this *Queue) Stop() {
	// 等待数据完成
	if len(this.c) > 0 {
		time.Sleep(200 * time.Millisecond)
	}

	close(this.c)
	this.c = nil

	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}
}

// 添加索引
func (this *Queue) Index(index []string) {
	if len(index) == 0 {
		return
	}

	this.indexLocker.Lock()
	defer this.indexLocker.Unlock()

	if this.coll == nil {
		return
	}

	// 是否已存在
	for _, i := range this.indexes {
		if this.equalStrings(index, i) {
			return
		}
	}

	indexMap := map[string]bool{
		"item": true,
	}
	for _, i := range index {
		indexMap["params."+i] = true
	}
	err := this.coll.CreateIndex(indexMap)
	if err != nil {
		logs.Error(err)
	}

	this.indexes = append(this.indexes, index)
}

// 增加值
// 只支持int, int32, int64, float32, float64
func (this *Queue) increase(value maps.Map, inc maps.Map) maps.Map {
	if inc == nil {
		return maps.Map{}
	}
	if value == nil {
		return inc
	}
	for k, v := range inc {
		v1, ok := value[k]
		if !ok {
			value[k] = v
			continue
		}
		switch v2 := v1.(type) {
		case int:
			v1 = v2 + types.Int(v)
		case int32:
			v1 = v2 + types.Int32(v)
		case int64:
			v1 = v2 + types.Int64(v)
		case float32:
			v1 = v2 + types.Float32(v)
		case float64:
			v1 = v2 + types.Float64(v)
		default:
			logs.Println("[teastats]queue increase not match:", reflect.TypeOf(v1).Kind())
		}
		value[k] = v1
	}

	return value
}

// 对比字符串数组看是否相等
func (this *Queue) equalStrings(strings1 []string, strings2 []string) bool {
	for _, s1 := range strings1 {
		if !lists.ContainsString(strings2, s1) {
			return false
		}
	}
	for _, s2 := range strings2 {
		if !lists.ContainsString(strings1, s2) {
			return false
		}
	}
	return true
}
