package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 请求数统计
type RequestAllPeriodFilter struct {
	CounterFilter
}

func (this *RequestAllPeriodFilter) Name() string {
	return "请求数统计"
}

func (this *RequestAllPeriodFilter) Codes() []string {
	return []string{
		"request.all.second",
		"request.all.minute",
		"request.all.hour",
		"request.all.day",
		"request.all.week",
		"request.all.month",
		"request.all.year",
	}
}

func (this *RequestAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *RequestAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *RequestAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	this.ApplyFilter(accessLog, nil, maps.Map{
		"count": 1,
	})
}

func (this *RequestAllPeriodFilter) Stop() {
	this.StopFilter()
}
