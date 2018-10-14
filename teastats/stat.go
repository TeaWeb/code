package teastats

import (
	"context"
	"fmt"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"sync"
	"time"
)

type IncrementOperation struct {
	coll   *teamongo.Collection
	filter map[string]interface{}
	init   map[string]interface{}
	field  string
}

func (this *IncrementOperation) uniqueId() string {
	keys := []string{}
	for key := range this.filter {
		keys = append(keys, key)
	}
	lists.Sort(keys, func(i int, j int) bool {
		return keys[i] < keys[j]
	})

	uniqueId := fmt.Sprintf("%p", this.coll)
	for _, key := range keys {
		uniqueId += "@" + types.String(this.filter[key])
	}
	return uniqueId + "@" + this.field
}

type Stat struct {
	operations []*IncrementOperation

	once   sync.Once
	locker sync.Mutex
}

func (this *Stat) initOnce() {
	this.once.Do(func() {
		timers.Loop(1*time.Second, func(looper *timers.Looper) {
			this.locker.Lock()
			defer this.locker.Unlock()

			lastUniqueId := ""
			count := 0
			end := 0

			for index, op := range this.operations {
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
				firstOP := this.operations[0]
				_, err := firstOP.coll.UpdateOne(context.Background(), firstOP.filter, map[string]interface{}{
					"$set": firstOP.init,
					"$inc": map[string]interface{}{
						firstOP.field: count,
					},
				}, updateopt.OptUpsert(true))
				if err != nil {
					logs.Error(err)
				}

				this.operations = this.operations[end+1:]
			}
		})
	})
}

func (this *Stat) Increase(collection *teamongo.Collection, filter map[string]interface{}, init map[string]interface{}, field string) {
	this.initOnce()

	if collection == nil {
		return
	}

	this.locker.Lock()
	defer this.locker.Unlock()
	this.operations = append(this.operations, &IncrementOperation{
		coll:   collection,
		filter: filter,
		init:   init,
		field:  field,
	})
}
