package teastats

import (
	"context"
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson"
	"strings"
	"time"
)

type DailyRequestsStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Day      string `bson:"day" json:"day"`           // 日期，格式为：Ymd
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *DailyRequestsStat) Init() {
	coll := findCollection("stats.requests.daily", nil)
	coll.CreateIndex(map[string]bool{
		"day": true,
	})
	coll.CreateIndex(map[string]bool{
		"day":      true,
		"serverId": true,
	})
}

func (this *DailyRequestsStat) Process(accessLog *tealogs.AccessLog) {
	day := timeutil.Format("Ymd")
	coll := findCollection("stats.requests.daily", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, "count")
}

func (this *DailyRequestsStat) ListLatestDays(serverId string, days int) []map[string]interface{} {
	if days <= 0 {
		days = 7
	}

	result := []map[string]interface{}{}
	for i := days - 1; i >= 0; i -- {
		day := timeutil.Format("Ymd", time.Now().AddDate(0, 0, -i))
		total := this.SumDayRequests(serverId, []string{day})
		result = append(result, map[string]interface{}{
			"day":   day,
			"total": total,
		})
	}
	return result
}

func (this *DailyRequestsStat) SumDayRequests(serverId string, days []string) int64 {
	if len(days) == 0 {
		return 0
	}
	sumColl := findCollection("stats.requests.daily", nil)

	pipelines := bson.NewArray()
	err := bson.UnmarshalExtJSON([]byte(`[
	{
		"$match": {
			"serverId": "` + serverId + `",
			"day": {
				"$in": [ "` + strings.Join(days, "\", \"") + `" ]
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
]`), true, &pipelines)
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
