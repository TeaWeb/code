package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// PV统计
type PVPagePeriodFilter struct {
	CounterFilter
}

func (this *PVPagePeriodFilter) Name() string {
	return "URL PV统计"
}

func (this *PVPagePeriodFilter) Codes() []string {
	return []string{
		"pv.page.second",
		"pv.page.minute",
		"pv.page.hour",
		"pv.page.day",
		"pv.page.week",
		"pv.page.month",
		"pv.page.year",
	}
}

func (this *PVPagePeriodFilter) Indexes() []string {
	return []string{"path"}
}

func (this *PVPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *PVPagePeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"path": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *PVPagePeriodFilter) Stop() {
	this.StopFilter()
}
