package teamongo

import (
	"context"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type Collection struct {
	*mongo.Collection
}

func FindCollection(collName string) *Collection {
	return &Collection{
		SharedClient().Database("teaweb").Collection(collName),
	}
}

// 创建索引
func (this *Collection) CreateIndex(indexes map[string]bool) error {
	manager := this.Indexes()

	doc := bson.NewDocument()

	// 对key进行排序
	keys := maps.NewMap(indexes).Keys()
	lists.Sort(keys, func(i int, j int) bool {
		return keys[i].(string) < keys[j].(string)
	})

	for _, key := range keys {
		index := key.(string)
		b := indexes[index]
		if b {
			doc.Append(bson.EC.Int32(index, 1))
		} else {
			doc.Append(bson.EC.Int32(index, -1))
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := manager.CreateOne(ctx, mongo.IndexModel{
		Keys:    doc,
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})
	return err
}
