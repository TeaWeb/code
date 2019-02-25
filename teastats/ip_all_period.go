package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// IP统计
type IPAllPeriodFilter struct {
	CounterFilter
}

func (this *IPAllPeriodFilter) Name() string {
	return "IP统计"
}

func (this *IPAllPeriodFilter) Codes() []string {
	return []string{
		"ip.all.second",
		"ip.all.minute",
		"ip.all.hour",
		"ip.all.day",
		"ip.all.week",
		"ip.all.month",
		"ip.all.year",
	}
}

func (this *IPAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *IPAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *IPAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if !this.CheckNewIP(accessLog, "") {
		return
	}

	this.ApplyFilter(accessLog, nil, maps.Map{
		"count": 1,
	})
}

func (this *IPAllPeriodFilter) Stop() {
	this.StopFilter()
}
