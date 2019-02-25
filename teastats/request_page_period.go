package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 请求数统计
type RequestPagePeriodFilter struct {
	CounterFilter
}

func (this *RequestPagePeriodFilter) Name() string {
	return "URL请求数统计"
}

func (this *RequestPagePeriodFilter) Codes() []string {
	return []string{
		"request.page.second",
		"request.page.minute",
		"request.page.hour",
		"request.page.day",
		"request.page.week",
		"request.page.month",
		"request.page.year",
	}
}

func (this *RequestPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *RequestPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *RequestPagePeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *RequestPagePeriodFilter) Stop() {
	this.StopFilter()
}
