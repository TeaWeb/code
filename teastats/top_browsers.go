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

type TopBrowserStat struct {
	Stat

	ServerId string  `bson:"serverId" json:"serverId"` // 服务ID
	Month    string  `bson:"month" json:"month"`       // 月份
	Family   string  `bson:"family" json:"family"`     // 浏览器
	Version  string  `bson:"version" json:"version"`   // 版本
	Count    int64   `bson:"count" json:"count"`       // 访问数量
	Percent  float64 `bson:"percent" json:"percent"`   // 比例
}

func (this *TopBrowserStat) Init() {
	coll := findCollection("stats.top.browsers.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"family":   true,
		"version":  true,
		"month":    true,
	})
	coll.CreateIndex(map[string]bool{
		"count": false,
	})
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
}

func (this *TopBrowserStat) Process(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Client.Browser.Family) == 0 {
		return
	}
	family := accessLog.Extend.Client.Browser.Family
	version := accessLog.Extend.Client.Browser.Major

	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.browsers.monthly", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"family":   family,
		"version":  version,
		"month":    month,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"family":   family,
		"version":  version,
		"month":    month,
	}, "count")
}

func (this *TopBrowserStat) List(serverId string, size int64) (result []TopBrowserStat) {
	if size <= 0 {
		size = 10
	}

	result = []TopBrowserStat{}

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
	coll := findCollection("stats.top.browsers.monthly", nil)
	pipelines, err := teamongo.BSONArrayBytes([]byte(`[
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
			"_id": {  "$concat" : [ "$family", " ", "$version" ]},
			"count": {
				"$sum": "$count"
			},
			"family": {
				"$first": "$family"
			},
			"version": {
				"$first": "$version"
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

	count := int64(0)
	for cursor.Next(context.Background()) {
		one := TopBrowserStat{}
		err := cursor.Decode(&one)
		if err == nil {
			if one.Family == "Other" || count >= size {
				continue
			}

			count ++

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

	return
}
