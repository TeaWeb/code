package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 浏览器统计
type BrowserAllPeriodFilter struct {
	CounterFilter
}

func (this *BrowserAllPeriodFilter) Name() string {
	return "浏览器统计"
}

func (this *BrowserAllPeriodFilter) Codes() []string {
	return []string{
		"browser.all.second",
		"browser.all.minute",
		"browser.all.hour",
		"browser.all.day",
		"browser.all.week",
		"browser.all.month",
		"browser.all.year",
	}
}

func (this *BrowserAllPeriodFilter) Indexes() []string {
	return []string{"family", "major"}
}

func (this *BrowserAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *BrowserAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Client.Browser.Family) == 0 {
		return
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV ++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Client.Browser.Family+"_"+accessLog.Extend.Client.Browser.Major) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Client.Browser.Family+"_"+accessLog.Extend.Client.Browser.Major) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"family": accessLog.Extend.Client.Browser.Family,
		"major":  accessLog.Extend.Client.Browser.Major,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *BrowserAllPeriodFilter) Stop() {
	this.StopFilter()
}
