package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// UV统计
type UVAllPeriodFilter struct {
	CounterFilter
}

func (this *UVAllPeriodFilter) Name() string {
	return "UV统计"
}

func (this *UVAllPeriodFilter) Codes() []string {
	return []string{
		"uv.all.second",
		"uv.all.minute",
		"uv.all.hour",
		"uv.all.day",
		"uv.all.week",
		"uv.all.month",
		"uv.all.year",
	}
}

func (this *UVAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *UVAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *UVAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if !this.CheckNewUV(accessLog, "") {
		return
	}

	this.ApplyFilter(accessLog, nil, maps.Map{
		"count": 1,
	})
}

func (this *UVAllPeriodFilter) Stop() {
	this.StopFilter()
}
