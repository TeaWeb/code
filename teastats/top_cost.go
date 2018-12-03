package teastats

import (
	"context"
	"fmt"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"strings"
	"time"
)

type TopCostStat struct {
	ServerId  string  `bson:"serverId" json:"serverId"`   // 服务ID
	Month     string  `bson:"month" json:"month"`         // 月份
	URL       string  `bson:"url" json:"url"`             // URL
	Cost      float64 `bson:"cost" json:"cost"`           // 平均耗时
	TotalCost float64 `bson:"totalCost" json:"totalCost"` // 总耗时
	Count     int64   `bson:"count" json:"count"`         // 请求数量
}

func (this *TopCostStat) Init() {
	coll := findCollection("stats.top.cost.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"month":    true,
		"url":      true,
	})
	coll.CreateIndex(map[string]bool{
		"cost": false,
	})
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
}

func (this *TopCostStat) Process(accessLog *tealogs.AccessLog) {
	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.cost.monthly", this.Init)

	url := accessLog.Scheme + "://" + accessLog.Host + accessLog.RequestURI

	filter := map[string]interface{}{
		"serverId": accessLog.ServerId,
		"url":      url,
		"month":    month,
	}

	stat := map[string]interface{}{
		"$set": map[string]interface{}{
			"serverId": accessLog.ServerId,
			"url":      url,
			"month":    month,
			"cost":     accessLog.RequestTime,
		},
		"$inc": map[string]interface{}{
			"count":     1,
			"totalCost": accessLog.RequestTime,
		},
	}

	result := coll.FindOneAndUpdate(context.Background(), filter, stat, findopt.OptUpsert(true), findopt.Projection(map[string]int{
		"_id":       1,
		"totalCost": 1,
		"count":     1,
	}))

	m := map[string]interface{}{}
	if result.Decode(m) != mongo.ErrNoDocuments {
		count := types.Int64(m["count"]) + 1
		totalCost := types.Float64(m["totalCost"]) + accessLog.RequestTime
		avgCost := totalCost / float64(count)
		_, err := coll.UpdateOne(context.Background(), bson.NewDocument(bson.EC.Interface("_id", m["_id"])), map[string]interface{}{
			"$set": map[string]interface{}{
				"cost": avgCost,
			},
		})
		if err != nil {
			logs.Error(err)
		}
	}
}

func (this *TopCostStat) List(serverId string, size int64) (result []TopCostStat) {
	if size <= 0 {
		size = 10
	}

	result = []TopCostStat{}

	// 最近两个月
	months := []string{}
	month1 := timeutil.Format("Ym")
	month2 := timeutil.Format("Ym", time.Now().AddDate(0, -1, 0))
	if month1 != month2 {
		months = append(months, month1, month2)
	} else {
		months = append(months, month1)
	}

	// 开始查找
	coll := findCollection("stats.top.cost.monthly", nil)
	pipelines, err := teamongo.JSONArrayBytes([]byte(`[
	{
		"$match": {
			"serverId": "` + serverId + `",
			"month": {
				"$in": [ "` + strings.Join(months, "\", \"") + `" ]
			}
		}
	},
	{
		"$group": {
			"_id": "$url",
			"url": {
				"$first": "$url"
			},
			"cost": {
				"$first": "$cost"
			}
		}
	},
	{
		"$sort": {
			"cost": -1
		}
	},
	{
		"$limit": ` + fmt.Sprintf("%d", size) + `
	}
]`))
	if err != nil {
		return
	}
	cursor, err := coll.Aggregate(context.Background(), pipelines)
	if err != nil {
		logs.Error(err)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		one := TopCostStat{}
		err := cursor.Decode(&one)
		if err == nil {
			result = append(result, one)
		} else {
			logs.Error(err)
		}
	}

	return
}
