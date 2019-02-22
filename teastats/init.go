package teastats

import (
	"github.com/TeaWeb/code/tealogs"
)

func init() {
	// 注册筛选器
	RegisterFilter(
		new(RequestAllPeriodFilter),
		new(RequestPagePeriodFilter),
		new(StatusAllPeriodFilter),
		new(StatusPagePeriodFilter),
		new(TrafficAllPeriodFilter),
		new(TrafficPagePeriodFilter),
		new(PVAllPeriodFilter),
		new(PVPagePeriodFilter),
		new(UVAllPeriodFilter),
		new(UVPagePeriodFilter),
		new(IPAllPeriodFilter),
		new(IPPagePeriodFilter),
		new(MethodAllPeriodFilter),
		new(MethodPagePeriodFilter),
		new(CostAllPeriodFilter),
		new(CostPagePeriodFilter),
		new(RefererDomainPeriodFilter),
		new(RefererURLPeriodFilter),
		new(LandingPagePeriodFilter),
		new(BackendAllPeriodFilter),
		new(LocationAllPeriodFilter),
		new(RewriteAllPeriodFilter),
		new(FastcgiAllPeriodFilter),
		new(DeviceAllPeriodFilter),
		new(OSAllPeriodFilter),
		new(BrowserAllPeriodFilter),
		new(RegionAllPeriodFilter),
		new(CityAllPeriodFilter),
	)

	// 注册AccessLogHook
	tealogs.AddAccessLogHook(&tealogs.AccessLogHook{
		Process: func(accessLog *tealogs.AccessLog) (goNext bool) {
			if !accessLog.ShouldStat() {
				return true
			}
			serverQueue := FindServerQueue(accessLog.ServerId)
			if serverQueue == nil {
				return true
			}
			serverQueue.Filter(accessLog)
			return true
		},
	})
}
