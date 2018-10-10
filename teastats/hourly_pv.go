package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/TeaWeb/code/tealogs"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"strings"
	"github.com/iwind/TeaGo/logs"
	"time"
	"github.com/iwind/TeaGo/types"
)

type HourlyPVStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Hour     string `bson:"hour" json:"hour"`         // 小时，格式为：YmdH
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *HourlyPVStat) Init() {
	coll := findCollection("stats.pv.hourly", nil)
	coll.CreateIndex(map[string]bool{
		"hour": true,
	})
	coll.CreateIndex(map[string]bool{
		"hour":     true,
		"serverId": true,
	})
}

func (this *HourlyPVStat) Process(accessLog *tealogs.AccessLog) {
	if !strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		return
	}

	hour := timeutil.Format("YmdH")
	coll := findCollection("stats.pv.hourly", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"hour":     hour,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"hour":     hour,
	}, "count")
}

func (this *HourlyPVStat) ListLatestHours(hours int) []map[string]interface{} {
	if hours <= 0 {
		hours = 24
	}

	result := []map[string]interface{}{}
	for i := hours - 1; i >= 0; i -- {
		hour := timeutil.Format("YmdH", time.Now().Add(time.Duration(-i)*time.Hour))
		total := this.SumHourPV([]string{hour})
		result = append(result, map[string]interface{}{
			"hour":  hour,
			"total": total,
		})
	}
	return result
}

func (this *HourlyPVStat) SumHourPV(hours []string) int64 {
	if len(hours) == 0 {
		return 0
	}
	sumColl := findCollection("stats.pv.hourly", nil)
	sumCursor, err := sumColl.Aggregate(context.Background(), bson.NewArray(bson.VC.DocumentFromElements(
		bson.EC.SubDocumentFromElements(
			"$match",
			bson.EC.Interface("hour", map[string]interface{}{
				"$in": hours,
			}),
		),
	), bson.VC.DocumentFromElements(bson.EC.SubDocumentFromElements(
		"$group",
		bson.EC.Interface("_id", nil),
		bson.EC.SubDocumentFromElements("total", bson.EC.String("$sum", "$count")),
	))))
	if err != nil {
		logs.Error(err)
		return 0
	}
	defer sumCursor.Close(context.Background())

	if sumCursor.Next(context.Background()) {
		sumMap := map[string]interface{}{}
		err = sumCursor.Decode(sumMap)
		if err == nil {
			return types.Int64(sumMap["total"])
		} else {
			logs.Error(err)
		}
	}

	return 0
}
