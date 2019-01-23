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

type TopRequestStat struct {
	Stat

	ServerId string  `bson:"serverId" json:"serverId"` // 服务ID
	Month    string  `bson:"month" json:"month"`       // 月份
	URL      string  `bson:"url" json:"url"`           // URL
	Count    int64   `bson:"count" json:"count"`       // 耗时
	Percent  float64 `bson:"percent" json:"percent"`   // 占比
}

func (this *TopRequestStat) Init() {
	coll := findCollection("stats.top.requests.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"month":    true,
		"url":      true,
	})
	coll.CreateIndex(map[string]bool{
		"count": false,
	})
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
}

func (this *TopRequestStat) Process(accessLog *tealogs.AccessLog) {
	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.requests.monthly", this.Init)

	url := accessLog.Scheme + "://" + accessLog.Host + accessLog.RequestURI

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"url":      url,
		"month":    month,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"url":      url,
		"month":    month,
	}, "count")
}

// 列出所有的排名
func (this *TopRequestStat) List(serverId string, size int64) (result []TopRequestStat) {
	if size <= 0 {
		size = 10
	}

	result = []TopRequestStat{}

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
	coll := findCollection("stats.top.requests.monthly", nil)
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
			"_id": "$url",
			"count": {
				"$sum": "$count"
			},
			"url": {
				"$first": "$url"
			}
		}
	},
	{
		"$sort": {
			"count": -1
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
		one := TopRequestStat{}
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

	return
}
