package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 设备统计
type DeviceAllPeriodFilter struct {
	CounterFilter
}

func (this *DeviceAllPeriodFilter) Name() string {
	return "设备统计"
}

func (this *DeviceAllPeriodFilter) Codes() []string {
	return []string{
		"device.all.second",
		"device.all.minute",
		"device.all.hour",
		"device.all.day",
		"device.all.week",
		"device.all.month",
		"device.all.year",
	}
}

func (this *DeviceAllPeriodFilter) Indexes() []string {
	return []string{"family", "model"}
}

func (this *DeviceAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *DeviceAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.Extend.Client.Device.Family) == 0 {
		return
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV ++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Client.Device.Family+"_"+accessLog.Extend.Client.Device.Model) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Client.Device.Family+"_"+accessLog.Extend.Client.Device.Model) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"family": accessLog.Extend.Client.Device.Family,
		"model":  accessLog.Extend.Client.Device.Model,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *DeviceAllPeriodFilter) Stop() {
	this.StopFilter()
}
