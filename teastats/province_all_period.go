package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 省份统计
type ProvinceAllPeriodFilter struct {
	CounterFilter
}

func (this *ProvinceAllPeriodFilter) Name() string {
	return "省份统计"
}

func (this *ProvinceAllPeriodFilter) Codes() []string {
	return []string{
		"province.all.second",
		"province.all.minute",
		"province.all.hour",
		"province.all.day",
		"province.all.week",
		"province.all.month",
		"province.all.year",
	}
}

func (this *ProvinceAllPeriodFilter) Indexes() []string {
	return []string{"region", "province"}
}

func (this *ProvinceAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *ProvinceAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Geo.State) == 0 {
		return
	}

	// 中国特区
	if accessLog.Extend.Geo.Region == "台湾" {
		accessLog.Extend.Geo.Region = "中国台湾"
	} else if accessLog.Extend.Geo.Region == "香港" {
		accessLog.Extend.Geo.Region = "中国香港"
	} else if accessLog.Extend.Geo.Region == "澳门" {
		accessLog.Extend.Geo.Region = "中国澳门"
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV ++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Geo.Region+accessLog.Extend.Geo.State) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Geo.Region+accessLog.Extend.Geo.State) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"region":   accessLog.Extend.Geo.Region,
		"province": accessLog.Extend.Geo.State,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *ProvinceAllPeriodFilter) Stop() {
	this.StopFilter()
}
