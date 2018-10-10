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

type MonthlyPVStat struct {
	Stat

	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份，格式为：Ym
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *MonthlyPVStat) Init() {
	coll := findCollection("stats.pv.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
	coll.CreateIndex(map[string]bool{
		"month":    true,
		"serverId": true,
	})
}

func (this *MonthlyPVStat) Process(accessLog *tealogs.AccessLog) {
	if !strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		return
	}

	month := timeutil.Format("Ym")
	coll := findCollection("stats.pv.monthly", this.Init)

	this.Increase(coll, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"month":    month,
	}, map[string]interface{}{
		"serverId": accessLog.ServerId,
		"month":    month,
	}, "count")
}

func (this *MonthlyPVStat) ListLatestMonths(months int) []map[string]interface{} {
	if months <= 0 {
		months = 12
	}

	result := []map[string]interface{}{}
	for i := months - 1; i >= 0; i -- {
		month := timeutil.Format("Ym", time.Now().AddDate(0, -i, 0))
		total := this.SumMonthPV([]string{month})
		result = append(result, map[string]interface{}{
			"month": month,
			"total": total,
		})
	}
	return result
}

func (this *MonthlyPVStat) SumMonthPV(months []string) int64 {
	if len(months) == 0 {
		return 0
	}
	sumColl := findCollection("stats.pv.monthly", nil)
	sumCursor, err := sumColl.Aggregate(context.Background(), bson.NewArray(bson.VC.DocumentFromElements(
		bson.EC.SubDocumentFromElements(
			"$match",
			bson.EC.Interface("month", map[string]interface{}{
				"$in": months,
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
