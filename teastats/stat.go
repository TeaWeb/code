package teastats

import (
	"context"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"sync"
	"time"
)

// 统计基本定义
type Stat struct {
	incOperations []*IncrementOperation
	avgOperations []*AvgOperation

	once   sync.Once
	locker sync.Mutex
}

// 初始化
func (this *Stat) initOnce() {
	this.once.Do(func() {
		timers.Loop(1*time.Second, func(looper *timers.Looper) {
			this.locker.Lock()
			defer this.locker.Unlock()

			// 增加操作
			{
				lastUniqueId := ""
				count := 0
				end := 0

				for index, op := range this.incOperations {
					uniqueId := op.uniqueId()
					if len(lastUniqueId) == 0 {
						lastUniqueId = uniqueId
						count ++
						end = index
						continue
					}

					if lastUniqueId == uniqueId {
						count ++
						end = index
						continue
					}

					// 不相同，则终止
					break
				}

				if count > 0 {
					firstOP := this.incOperations[0]
					_, err := firstOP.coll.UpdateOne(context.Background(), firstOP.filter, map[string]interface{}{
						"$set": firstOP.init,
						"$inc": map[string]interface{}{
							firstOP.field: count,
						},
					}, updateopt.OptUpsert(true))
					if err != nil {
						logs.Error(err)
					}

					this.incOperations = this.incOperations[end+1:]
				}
			}

			// 平均值操作
			{
				dataMap := map[string]maps.Map{} // uniqueId => { count, sum, op }
				for _, op := range this.avgOperations {
					uniqueId := op.uniqueId()
					m, found := dataMap[uniqueId]
					if found {
						m["count"] = m.GetInt("count") + op.count
						m["sum"] = m.GetFloat64("sum") + op.sum
					} else {
						m = maps.Map{}
						m["count"] = op.count
						m["sum"] = op.sum
						m["op"] = op
						dataMap[uniqueId] = m
					}
				}
				for _, m := range dataMap {
					op := m["op"].(*AvgOperation)
					count := m.GetInt("count")
					sum := m.GetFloat64("sum")
					_, err := op.coll.UpdateOne(context.Background(), op.filter, map[string]interface{}{
						"$set": op.init,
						"$inc": map[string]interface{}{
							op.countField: count,
							op.sumField:   sum,
						},
					}, updateopt.OptUpsert(true))
					if err != nil {
						logs.Error(err)
					}
				}
				this.avgOperations = []*AvgOperation{}
			}
		})
	})
}

// 为某个字段做增加操作
func (this *Stat) Increase(collection *teamongo.Collection, filter map[string]interface{}, init map[string]interface{}, field string) {
	this.initOnce()

	if collection == nil {
		return
	}

	this.locker.Lock()
	defer this.locker.Unlock()
	this.incOperations = append(this.incOperations, &IncrementOperation{
		coll:   collection,
		filter: filter,
		init:   init,
		field:  field,
	})
}

// 为某个字段做平均值计算操作
func (this *Stat) Avg(collection *teamongo.Collection, filter map[string]interface{}, init map[string]interface{}, countField string, count int, sumField string, sum float64) {
	this.initOnce()

	if collection == nil {
		return
	}

	this.locker.Lock()
	defer this.locker.Unlock()
	this.avgOperations = append(this.avgOperations, &AvgOperation{
		coll:       collection,
		filter:     filter,
		init:       init,
		countField: countField,
		count:      count,
		sumField:   sumField,
		sum:        sum,
	})
}
