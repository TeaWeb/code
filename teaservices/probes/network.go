package probes

import (
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/logs"
	"github.com/shirou/gopsutil/net"
	"math"
	"time"
)

type NetworkProbe struct {
	Probe

	lastBytesSent     uint64
	lastBytesReceived uint64
	lastTime          time.Time
}

func (this *NetworkProbe) Run() {
	this.InitOnce(func() {
		this.lastTime = time.Now()

		widget := teaplugins.NewWidget()
		widget.Dashboard = true
		widget.Group = teaplugins.WidgetGroupSystem
		widget.Name = "网络相关"

		{
			chart := teacharts.NewGaugeChart()
			chart.Name = "出口带宽"
			chart.Detail = "兆字节"
			chart.Unit = "MB"
			chart.Max = 10
			widget.AddChart(chart)
		}

		{
			chart := teacharts.NewGaugeChart()
			chart.Name = "入口带宽"
			chart.Detail = "兆字节"
			chart.Unit = "MB"
			chart.Max = 10
			widget.AddChart(chart)
		}

		this.Plugin.AddWidget(widget)

		widget.OnReload(func() {
			this.Run()
		})
	})

	stats, err := net.IOCounters(false)
	if err != nil {
		logs.Error(err)
	} else if len(stats) > 0 {
		countBytesSent := uint64(0)
		countBytesReceived := uint64(0)

		seconds := time.Since(this.lastTime).Seconds()
		this.lastTime = time.Now()

		for _, stat := range stats {
			countBytesSent += stat.BytesSent
			countBytesReceived += stat.BytesRecv
		}

		{
			chart := this.Plugin.Widgets[0].Charts[0].(*teacharts.GaugeChart)
			if this.lastBytesSent == 0 {
				chart.Value = 0
			} else {
				chart.Max = float64(((countBytesSent-this.lastBytesSent)/1024/1024/10 + 1) * 10)
				chart.Value = math.Ceil(float64(countBytesSent-this.lastBytesSent)*100/1024/1024/seconds) / 100
			}
			this.lastBytesSent = countBytesSent
		}
		{
			chart := this.Plugin.Widgets[0].Charts[1].(*teacharts.GaugeChart)
			if this.lastBytesReceived == 0 {
				chart.Value = 0
			} else {
				chart.Max = float64(((countBytesReceived-this.lastBytesReceived)/1024/1024/10 + 1) * 10)
				chart.Value = math.Ceil(float64(countBytesReceived-this.lastBytesReceived)*100/1024/1024/seconds) / 100
			}
			this.lastBytesReceived = countBytesReceived
		}
	}
}
