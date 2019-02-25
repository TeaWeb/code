package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 操作系统统计
type OSAllPeriodFilter struct {
	CounterFilter
}

func (this *OSAllPeriodFilter) Name() string {
	return "操作系统统计"
}

func (this *OSAllPeriodFilter) Codes() []string {
	return []string{
		"os.all.second",
		"os.all.minute",
		"os.all.hour",
		"os.all.day",
		"os.all.week",
		"os.all.month",
		"os.all.year",
	}
}

func (this *OSAllPeriodFilter) Indexes() []string {
	return []string{"family", "major"}
}

func (this *OSAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *OSAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Client.OS.Family) == 0 {
		return
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV ++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Client.OS.Family+"_"+accessLog.Extend.Client.OS.Major) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Client.OS.Family+"_"+accessLog.Extend.Client.OS.Major) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"family": accessLog.Extend.Client.OS.Family,
		"major":  accessLog.Extend.Client.OS.Major,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *OSAllPeriodFilter) Stop() {
	this.StopFilter()
}
