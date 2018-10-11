package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/TeaWeb/code/tealogs"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"strings"
	"github.com/iwind/TeaGo/logs"
	"time"
	"github.com/iwind/TeaGo/types"
)

type DailyUVStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Day      string `bson:"day" json:"day"`           // 日期，格式为：Ymd
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *DailyUVStat) Init() {
	coll := findCollection("stats.uv.daily", nil)
	coll.CreateIndex(map[string]bool{
		"day": true,
	})
	coll.CreateIndex(map[string]bool{
		"day":      true,
		"serverId": true,
	})
}

func (this *DailyUVStat) Process(accessLog *tealogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	day := timeutil.Format("Ymd")

	// 是否已存在
	result := findCollection("logs."+day, nil).FindOne(context.Background(), bson.NewDocument(bson.EC.String("remoteAddr", accessLog.RemoteAddr), bson.EC.String("serverId", accessLog.ServerId)), findopt.Projection(map[string]int{
		"id": 1,
	}))

	existAccessLog := map[string]interface{}{}
	if result.Decode(existAccessLog) != mongo.ErrNoDocuments {
		return
	}

	coll := findCollection("stats.uv.daily", this.Init)
	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, "count")
}

func (this *DailyUVStat) ListLatestDays(days int) []map[string]interface{} {
	if days <= 0 {
		days = 7
	}

	result := []map[string]interface{}{}
	for i := days - 1; i >= 0; i -- {
		day := timeutil.Format("Ymd", time.Now().AddDate(0, 0, -i))
		total := this.SumDayUV([]string{day})
		result = append(result, map[string]interface{}{
			"day":   day,
			"total": total,
		})
	}
	return result
}

func (this *DailyUVStat) SumDayUV(days []string) int64 {
	if len(days) == 0 {
		return 0
	}
	sumColl := findCollection("stats.uv.daily", nil)
	pipelines, err := bson.ParseExtJSONArray(`[
	{
		"$match": {
			"day": {
				"$in": [ "` + strings.Join(days, ", ") + `" ]
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
