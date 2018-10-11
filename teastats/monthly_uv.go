package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/TeaWeb/code/tealogs"
	"strings"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaGo/logs"
	"time"
	"github.com/iwind/TeaGo/types"
)

type MonthlyUVStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份，格式为：Ym
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *MonthlyUVStat) Init() {
	coll := findCollection("stats.uv.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
	coll.CreateIndex(map[string]bool{
		"month":    true,
		"serverId": true,
	})
}

func (this *MonthlyUVStat) Process(accessLog *tealogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	month := timeutil.Format("Ym")

	// 是否已存在
	result := findCollection("logs."+timeutil.Format("Ymd"), nil).FindOne(context.Background(), bson.NewDocument(bson.EC.String("remoteAddr", accessLog.RemoteAddr), bson.EC.String("serverId", accessLog.ServerId)), findopt.Projection(map[string]int{
		"id": 1,
	}))

	existAccessLog := map[string]interface{}{}
	if result.Decode(existAccessLog) != mongo.ErrNoDocuments {
		return
	}

	coll := findCollection("stats.uv.monthly", this.Init)
	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"month":    month,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"month":    month,
	}, "count")
}

func (this *MonthlyUVStat) ListLatestMonths(months int) []map[string]interface{} {
	if months <= 0 {
		months = 12
	}

	result := []map[string]interface{}{}
	for i := months - 1; i >= 0; i -- {
		month := timeutil.Format("Ym", time.Now().AddDate(0, -i, 0))
		total := this.SumMonthUV([]string{month})
		result = append(result, map[string]interface{}{
			"month": month,
			"total": total,
		})
	}
	return result
}

func (this *MonthlyUVStat) SumMonthUV(months []string) int64 {
	if len(months) == 0 {
		return 0
	}
	sumColl := findCollection("stats.uv.monthly", nil)
	pipelines, err := bson.ParseExtJSONArray(`[
	{
		"$match": {
			"month": {
				"$in": [ "` + strings.Join(months, "\", \"") + `" ]
			}
		}
	},
	{
		"$group": {
			"_id": null,
			"total": {
				"$sum": "$count"
			}
		}
	}
]`)
	if err != nil {
		logs.Error(err)
		return 0
	}

	sumCursor, err := sumColl.Aggregate(context.Background(), pipelines)
	if err != nil {
		logs.Error(err)
		return 0
	}
	defer sumCursor.Close(context.Background())

	if sumCursor.Next(context.Background()) {
		sumMap := map[string]interface{}{}
		err = sumCursor.Decode(&sumMap)
		if err == nil {
			return types.Int64(sumMap["total"])
		} else {
			logs.Error(err)
		}
	}

	return 0
}
