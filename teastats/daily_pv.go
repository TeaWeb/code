package teastats

import (
	"context"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"strings"
	"time"
)

type DailyPVStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Day      string `bson:"day" json:"day"`           // 日期，格式为：Ymd
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *DailyPVStat) Init() {
	coll := findCollection("stats.pv.daily", nil)
	coll.CreateIndex(map[string]bool{
		"day": true,
	})
	coll.CreateIndex(map[string]bool{
		"day":      true,
		"serverId": true,
	})
}

func (this *DailyPVStat) Process(accessLog *tealogs.AccessLog) {
	if !strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		return
	}

	day := timeutil.Format("Ymd")
	coll := findCollection("stats.pv.daily", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, "count")
}

func (this *DailyPVStat) ListLatestDays(serverId string, days int) []map[string]interface{} {
	if days <= 0 {
		days = 7
	}

	result := []map[string]interface{}{}
	for i := days - 1; i >= 0; i -- {
		day := timeutil.Format("Ymd", time.Now().AddDate(0, 0, -i))
		total := this.SumDayPV(serverId, []string{day})
		result = append(result, map[string]interface{}{
			"day":   day,
			"total": total,
		})
	}
	return result
}

func (this *DailyPVStat) SumDayPV(serverId string, days []string) int64 {
	if len(days) == 0 {
		return 0
	}
	sumColl := findCollection("stats.pv.daily", nil)

	pipelines, err := teamongo.JSONArrayBytes([]byte(`[
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
]`))
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
