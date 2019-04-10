package teamongo

import (
	"context"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Collection struct {
	*mongo.Collection
}

func FindCollection(collName string) *Collection {
	return &Collection{
		SharedClient().Database(DatabaseName).Collection(collName),
	}
}

// 创建索引
func (this *Collection) CreateIndex(indexes map[string]bool) error {
	indexView := this.Indexes()

	doc := map[string]interface{}{}

	// 对key进行排序
	keys := maps.NewMap(indexes).Keys()
	lists.Sort(keys, func(i int, j int) bool {
		return keys[i].(string) < keys[j].(string)
	})

	for _, key := range keys {
		index := key.(string)
		b := indexes[index]
		if b {
			doc[index] = 1
		} else {
			doc[index] = -1
		}
	}

	// 检查是否已经存在
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := indexView.List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
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

	// 创建新的
	_, err = indexView.CreateOne(ctx, mongo.IndexModel{
		Keys:    doc,
		Options: options.Index().SetBackground(true),
	})
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
