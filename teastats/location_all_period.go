package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 路径规则请求统计
type LocationAllPeriodFilter struct {
	CounterFilter
}

func (this *LocationAllPeriodFilter) Name() string {
	return "路径规则请求统计"
}

func (this *LocationAllPeriodFilter) Codes() []string {
	return []string{
		"location.all.second",
		"location.all.minute",
		"location.all.hour",
		"location.all.day",
		"location.all.week",
		"location.all.month",
		"location.all.year",
	}
}

func (this *LocationAllPeriodFilter) Indexes() []string {
	return []string{"location"}
}

func (this *LocationAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *LocationAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.LocationId) == 0 {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"location": accessLog.LocationId,
	}, maps.Map{
		"count": 1,
	})
}

func (this *LocationAllPeriodFilter) Stop() {
	this.StopFilter()
}
