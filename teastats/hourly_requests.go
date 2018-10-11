package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/TeaWeb/code/tealogs"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"time"
	"strings"
)

type HourlyRequestsStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Hour     string `bson:"hour" json:"hour"`         // 小时，格式为：YmdH
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *HourlyRequestsStat) Init() {
	coll := findCollection("stats.requests.hourly", nil)
	coll.CreateIndex(map[string]bool{
		"hour": true,
	})
	coll.CreateIndex(map[string]bool{
		"hour":     true,
		"serverId": true,
	})
}

func (this *HourlyRequestsStat) Process(accessLog *tealogs.AccessLog) {
	hour := timeutil.Format("YmdH")
	coll := findCollection("stats.requests.hourly", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"hour":     hour,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"hour":     hour,
	}, "count")
}

func (this *HourlyRequestsStat) ListLatestHours(hours int) []map[string]interface{} {
	if hours <= 0 {
		hours = 24
	}

	result := []map[string]interface{}{}
	for i := hours - 1; i >= 0; i -- {
		hour := timeutil.Format("YmdH", time.Now().Add(time.Duration(-i)*time.Hour))
		total := this.SumHourRequests([]string{hour})
		result = append(result, map[string]interface{}{
			"hour":  hour,
			"total": total,
		})
	}
	return result
}

func (this *HourlyRequestsStat) SumHourRequests(hours []string) int64 {
	if len(hours) == 0 {
		return 0
	}
	sumColl := findCollection("stats.requests.hourly", nil)
	pipelines, err := bson.ParseExtJSONArray(`[
	{
		"$match": {
			"hour": {
				"$in": [ "` + strings.Join(hours, ", ") + `" ]
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
