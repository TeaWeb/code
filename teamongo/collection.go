package teamongo

import (
	"context"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"sync"
	"time"
)

var collMap = map[string]*Collection{}
var collLocker = sync.RWMutex{}

// 集合定义
type Collection struct {
	*mongo.Collection
}

// 取得共享的集合
func SharedCollection(collName string) *Collection {
	collLocker.RLock()
	coll, ok := collMap[collName]
	if ok {
		collLocker.RUnlock()
		return coll
	}
	collLocker.RUnlock()

	coll = &Collection{
		SharedClient().Database(DatabaseName).Collection(collName),
	}

	collLocker.Lock()
	collMap[collName] = coll
	collLocker.Unlock()

	return coll
}

// 创建索引
func (this *Collection) CreateIndex(fields ...*shared.IndexField) error {
	indexView := this.Indexes()

	doc := map[string]interface{}{}

	for _, field := range fields {
		if field.Asc {
			doc[field.Name] = 1
		} else {
			doc[field.Name] = -1
		}
	}

	// 检查是否已经存在
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := indexView.List(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			logs.Error(err)
		}
	}()
	for cursor.Next(ctx) {
		m := map[string]interface{}{}
		err = cursor.Decode(&m)
		if err != nil {
			return err
		}
		key, ok := m["key"]
		if !ok {
			continue
		}
		keyMap, ok := key.(map[string]interface{})
		if !ok {
			continue
		}
		if checkIndexEqual(doc, keyMap) {
			return nil
		}
	}

	bsonDoc := bsonx.Doc{}
	for _, field := range fields {
		if field.Asc {
			bsonDoc = bsonDoc.Append(field.Name, bsonx.Int32(1))
		} else {
			bsonDoc = bsonDoc.Append(field.Name, bsonx.Int32(-1))
		}
	}

	// 创建新的
	_, err = indexView.CreateOne(ctx, mongo.IndexModel{
		Keys:    bsonDoc,
		Options: options.Index().SetBackground(true),
	})
	return err
}

// 创建一组索引
func (this *Collection) CreateIndexes(fields ...[]*shared.IndexField) error {
	var err error = nil
	for _, f := range fields {
		err1 := this.CreateIndex(f...)
		if err1 != nil {
			err = err1
		}
	}
	return err
}

func checkIndexEqual(index1 map[string]interface{}, index2 map[string]interface{}) bool {
	if len(index1) != len(index2) {
		return false
	}
	for k, v := range index1 {
		v2, ok := index2[k]
		if !ok {
			return false
		}
		if types.Int(v) != types.Int(v2) {
			return false
		}
	}
	return true
}
