package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 请求方法统计
type MethodPagePeriodFilter struct {
	CounterFilter
}

func (this *MethodPagePeriodFilter) Name() string {
	return "URL请求方法统计"
}

func (this *MethodPagePeriodFilter) Codes() []string {
	return []string{
		"method.page.second",
		"method.page.minute",
		"method.page.hour",
		"method.page.day",
		"method.page.week",
		"method.page.month",
		"method.page.year",
	}
}

func (this *MethodPagePeriodFilter) Indexes() []string {
	return []string{"method", "page"}
}

func (this *MethodPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *MethodPagePeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"method": accessLog.RequestMethod,
		"page":   accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *MethodPagePeriodFilter) Stop() {
	this.StopFilter()
}
