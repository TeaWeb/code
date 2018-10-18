package probes

import (
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/shirou/gopsutil/disk"
	"strings"
	"time"
)

type DiskProbe struct {
	Probe
}

func (this *DiskProbe) Run() {
	this.InitOnce(func() {
		logs.Println("probe disk")

		widget := teaplugins.NewWidget()
		widget.Name = "文件系统"
		widget.Group = teaplugins.WidgetGroupSystem
		widget.Dashboard = true
		widget.OnForceReload(func() {
			this.Run()
		})

		t := time.Now()
		widget.OnReload(func() {
			if time.Since(t).Seconds() < 5 {
				return
			}

			t = time.Now()

			this.Run()
		})
		this.Plugin.AddWidget(widget)
	})

	widget := this.Plugin.Widgets[0]
	widget.Charts = []teacharts.ChartInterface{}

	partitions, err := disk.Partitions(true)
	if err != nil {
		logs.Error(err)
	} else {
		lists.Sort(partitions, func(i int, j int) bool {
			p1 := partitions[i]
			p2 := partitions[j]
			return p1.Mountpoint < p2.Mountpoint
		})

		for _, partition := range partitions {
			if !strings.Contains(partition.Device, "/") && !strings.Contains(partition.Device, "\\") {
				continue
			}

			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				logs.Error(err)
			} else {
				if usage.Total == 0 {
					continue
				}

				chart := teacharts.NewProgressBar()
				chart.SetUniqueId(stringutil.Md5(partition.Mountpoint))
				chart.Name = partition.Mountpoint

				if usage.Total == 0 {
					chart.Value = 0
				} else {
					chart.Value = usage.UsedPercent
					if usage.Total-usage.Used < 1024*1024*1024 { // 小于1G
						chart.Color = teacharts.ColorRed
					} else if usage.UsedPercent > 80 {
						chart.Color = teacharts.ColorRed
					} else if usage.UsedPercent > 40 {
						chart.Color = teacharts.ColorBlue
					} else {
						chart.Color = teacharts.ColorGreen
					}
				}

				if usage.Used == 0 && usage.Total == 0 {
					chart.Detail = ""
				} else {
					chart.Detail = formatBytes(usage.Used) + "/" + formatBytes(usage.Total)
				}

				widget.AddChart(chart)
			}
		}
	}
}
