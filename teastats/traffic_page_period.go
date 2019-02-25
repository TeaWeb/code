package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 流量统计
type TrafficPagePeriodFilter struct {
	CounterFilter
}

func (this *TrafficPagePeriodFilter) Name() string {
	return "URL流量统计"
}

func (this *TrafficPagePeriodFilter) Codes() []string {
	return []string{
		"traffic.page.second",
		"traffic.page.minute",
		"traffic.page.hour",
		"traffic.page.day",
		"traffic.page.week",
		"traffic.page.month",
		"traffic.page.year",
	}
}

func (this *TrafficPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *TrafficPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *TrafficPagePeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"bytes": accessLog.BytesSent,
	})
}

func (this *TrafficPagePeriodFilter) Stop() {
	this.StopFilter()
}
