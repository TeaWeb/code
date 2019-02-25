package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 流量统计
type TrafficAllPeriodFilter struct {
	CounterFilter
}

func (this *TrafficAllPeriodFilter) Name() string {
	return "流量统计"
}

func (this *TrafficAllPeriodFilter) Codes() []string {
	return []string{
		"traffic.all.second",
		"traffic.all.minute",
		"traffic.all.hour",
		"traffic.all.day",
		"traffic.all.week",
		"traffic.all.month",
		"traffic.all.year",
	}
}

func (this *TrafficAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *TrafficAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *TrafficAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	this.ApplyFilter(accessLog, nil, maps.Map{
		"bytes": accessLog.BytesSent,
	})
}

func (this *TrafficAllPeriodFilter) Stop() {
	this.StopFilter()
}
