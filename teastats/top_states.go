package teastats

import (
	"context"
	"fmt"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/time"
	"strings"
	"time"
)

type TopStateStat struct {
	Stat

	ServerId string  `bson:"serverId" json:"serverId"` // 服务ID
	Month    string  `bson:"month" json:"month"`       // 月份
	Region   string  `bson:"region" json:"region"`     // 地区
	State    string  `bson:"state" json:"state"`       // 省份|州
	Count    int64   `bson:"count" json:"count"`       // 访问数量
	Percent  float64 `bson:"percent" json:"percent"`   // 比例
}

func (this *TopStateStat) Init() {
	coll := findCollection("stats.top.states.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"state":    true,
		"region":   true,
		"month":    true,
	})
	coll.CreateIndex(map[string]bool{
		"count": false,
	})
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
}

func (this *TopStateStat) Process(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Geo.Region) == 0 || len(accessLog.Extend.Geo.State) == 0 {
		return
	}
	region := accessLog.Extend.Geo.Region
	state := accessLog.Extend.Geo.State

	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.states.monthly", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"region":   region,
		"state":    state,
		"month":    month,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"region":   region,
		"state":    state,
		"month":    month,
	}, "count")
}

func (this *TopStateStat) List(serverId string, size int64) (result []TopStateStat) {
	if size <= 0 {
		size = 10
	}

	result = []TopStateStat{}

	// 最近两个月
	months := []string{}
	month1 := timeutil.Format("Ym")
	month2 := timeutil.Format("Ym", time.Now().AddDate(0, -1, 0))
	if month1 != month2 {
		months = append(months, month1, month2)
	} else {
		months = append(months, month1)
	}

	// 总请求数量
	totalRequests := new(MonthlyRequestsStat).SumMonthRequests(serverId, months)

	// 开始查找
	coll := findCollection("stats.top.states.monthly", nil)
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
			"_id": {  "$concat" : [ "$region", " ", "$state" ]},
			"count": {
				"$sum": "$count"
			},
			"region": {
				"$first": "$region"
			},
			"state": {
				"$first": "$state"
			}
		}
	},
	{
		"$sort": {
			"count": -1
		}
	},
	{
		"$limit": ` + fmt.Sprintf("%d", size+1) + `
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
		one := TopStateStat{}
		err := cursor.Decode(&one)
		if err == nil {
			if totalRequests > 0 {
				one.Percent = float64(one.Count) / float64(totalRequests)
			} else {
				one.Percent = 0
			}

			result = append(result, one)
		} else {
			logs.Error(err)
		}
	}

	if len(result) > int(size) {
		result = result[:size]
	}

	return
}
