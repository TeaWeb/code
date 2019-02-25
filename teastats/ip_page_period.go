package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// IP统计
type IPPagePeriodFilter struct {
	CounterFilter
}

func (this *IPPagePeriodFilter) Name() string {
	return "URL IP统计"
}

func (this *IPPagePeriodFilter) Codes() []string {
	return []string{
		"ip.page.second",
		"ip.page.minute",
		"ip.page.hour",
		"ip.page.day",
		"ip.page.week",
		"ip.page.month",
		"ip.page.year",
	}
}

func (this *IPPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *IPPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *IPPagePeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if !this.CheckNewIP(accessLog, accessLog.RequestPath) {
		return
	}

	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *IPPagePeriodFilter) Stop() {
	this.StopFilter()
}
