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

type APIDailyRequestsStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	API      string `bson:"api" json:"api"`           // API path
	Day      string `bson:"day" json:"day"`           // 日期，格式为：Ymd
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *APIDailyRequestsStat) Init() {
	coll := findCollection("stats.api.requests.daily", nil)
	coll.CreateIndex(map[string]bool{
		"day": true,
	})
	coll.CreateIndex(map[string]bool{
		"day":      true,
		"serverId": true,
	})
	coll.CreateIndex(map[string]bool{
		"day":      true,
		"api":      true,
		"serverId": true,
	})
}

func (this *APIDailyRequestsStat) Process(accessLog *tealogs.AccessLog) {
	if len(accessLog.APIPath) == 0 {
		return
	}

	day := timeutil.Format("Ymd")
	coll := findCollection("stats.api.requests.daily", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"day":      day,
	}, "count")

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"api":      accessLog.APIPath,
		"day":      day,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"api":      accessLog.APIPath,
		"day":      day,
	}, "count")
}

func (this *APIDailyRequestsStat) ListLatestDays(serverId string, days int) []map[string]interface{} {
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

func (this *APIDailyRequestsStat) SumDayRequests(serverId string, days []string) int64 {
	if len(days) == 0 {
		return 0
	}
	sumColl := findCollection("stats.api.requests.daily", nil)

	pipelines, err := teamongo.JSONArrayBytes([]byte(`[
	{
		"$match": {
			"serverId": "` + serverId + `",
			"day": {
				"$in": [ "` + strings.Join(days, "\", \"") + `" ]
			},
			"api": null
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

func (this *APIDailyRequestsStat) ListLatestDaysForAPI(serverId string, apiPath string, days int) []map[string]interface{} {
	if days <= 0 {
		days = 7
	}

	result := []map[string]interface{}{}
	for i := days - 1; i >= 0; i -- {
		day := timeutil.Format("Ymd", time.Now().AddDate(0, 0, -i))
		total := this.SumDayRequestsForAPI(serverId, apiPath, []string{day})
		result = append(result, map[string]interface{}{
			"day":   day,
			"total": total,
		})
	}
	return result
}

func (this *APIDailyRequestsStat) SumDayRequestsForAPI(serverId string, apiPath string, days []string) int64 {
	if len(days) == 0 {
		return 0
	}
	sumColl := findCollection("stats.api.requests.daily", nil)

	pipelines, err := teamongo.JSONArrayBytes([]byte(`[
	{
		"$match": {
			"serverId": "` + serverId + `",
			"day": {
				"$in": [ "` + strings.Join(days, "\", \"") + `" ]
			},
			"api": "` + apiPath + `"
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
