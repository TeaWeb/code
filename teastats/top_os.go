package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/utils/time"
	"context"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/iwind/TeaGo/logs"
)

type TopOSStat struct {
	Stat

	ServerId string  `bson:"serverId" json:"serverId"` // 服务ID
	Month    string  `bson:"month" json:"month"`       // 月份
	Family   string  `bson:"family" json:"family"`     // 操作系统
	Version  string  `bson:"version" json:"version"`   // 版本
	Count    int64   `bson:"count" json:"count"`       // 访问数量
	Percent  float64 `bson:"percent" json:"percent"`   // 比例
}

func (this *TopOSStat) Init() {
	coll := findCollection("stats.top.os.monthly", nil)
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

func (this *TopOSStat) Process(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Client.OS.Family) == 0 {
		return
	}
	family := accessLog.Extend.Client.OS.Family
	version := accessLog.Extend.Client.OS.Major

	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.os.monthly", this.Init)

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

func (this *TopOSStat) List(size int64) (result []TopOSStat) {
	if size <= 0 {
		size = 10
	}

	result = []TopOSStat{}

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
	totalRequests := new(MonthlyRequestsStat).SumMonthRequests(months)

	// 开始查找
	coll := findCollection("stats.top.os.monthly", nil)
	cursor, err := coll.Find(context.Background(), map[string]interface{}{
		"month": map[string]interface{}{
			"$in": months,
		},
	}, findopt.Sort(map[string]interface{}{
		"count": -1,
	}), findopt.Limit(size))
	if err != nil {
		logs.Error(err)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		one := TopOSStat{}
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
