package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 城市统计
type CityAllPeriodFilter struct {
	CounterFilter
}

func (this *CityAllPeriodFilter) Name() string {
	return "城市统计"
}

func (this *CityAllPeriodFilter) Codes() []string {
	return []string{
		"city.all.second",
		"city.all.minute",
		"city.all.hour",
		"city.all.day",
		"city.all.week",
		"city.all.month",
		"city.all.year",
	}
}

func (this *CityAllPeriodFilter) Indexes() []string {
	return []string{"region", "province", "city"}
}

func (this *CityAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *CityAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Geo.City) == 0 {
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

	if this.CheckNewUV(accessLog, accessLog.Extend.Geo.Region+accessLog.Extend.Geo.State+accessLog.Extend.Geo.City) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Geo.Region+accessLog.Extend.Geo.State+accessLog.Extend.Geo.City) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"region":   accessLog.Extend.Geo.Region,
		"province": accessLog.Extend.Geo.State,
		"city":     accessLog.Extend.Geo.City,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *CityAllPeriodFilter) Stop() {
	this.StopFilter()
}
