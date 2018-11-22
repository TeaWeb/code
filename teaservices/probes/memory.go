package probes

import (
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/logs"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type MemoryProbe struct {
	Probe

	usageSizes []float64 // 单位为G
}

func (this *MemoryProbe) Run() {
	this.InitOnce(func() {
		logs.Println("probe memory")

		widget := teaplugins.NewWidget()
		widget.Name = "内存使用情况"
		widget.Dashboard = true
		widget.Group = teaplugins.WidgetGroupSystem

		{
			chart := teacharts.NewLineChart()
			widget.AddChart(chart)
		}

		{
			chart := teacharts.NewProgressBar()
			chart.Name = "物理内存"
			chart.Value = 0
			widget.AddChart(chart)
		}

		{
			chart := teacharts.NewProgressBar()
			chart.Name = "虚拟内存"
			chart.Value = 0
			widget.AddChart(chart)
		}

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

	{
		chart := this.Plugin.Widgets[0].Charts[1].(*teacharts.ProgressBar)
		stat, err := mem.VirtualMemory()

		if err != nil {
			logs.Error(err)
		} else {
			chart.Value = stat.UsedPercent
			if chart.Value < 50 {
				chart.Color = teacharts.ColorGreen
			} else if chart.Value > 80 {
				chart.Color = teacharts.ColorRed
			}
			chart.Detail = formatBytes(stat.Used) + "/" + formatBytes(stat.Total)

			// 线图
			{
				this.usageSizes = append(this.usageSizes, float64(stat.Used)/1024/1024/1024)
				if len(this.usageSizes) > 20 {
					this.usageSizes = this.usageSizes[len(this.usageSizes)-20:]
				}

				values := []interface{}{}
				labels := []string{}
				for _, size := range this.usageSizes {
					values = append(values, size)
					labels = append(labels, "")
				}

				chart := this.Plugin.Widgets[0].Charts[0].(*teacharts.LineChart)
				chart.Max = float64(stat.Total) / 1024 / 1024 / 1024
				chart.ResetLines()
				chart.YTickCount = 1

				color := teacharts.ColorBlue
				if stat.UsedPercent < 50 {
					color = teacharts.ColorGreen
				} else if stat.UsedPercent > 80 {
					color = teacharts.ColorRed
				}

				chart.AddLine(&teacharts.Line{
					Name:      "物理内存使用",
					Values:    values,
					Filled:    true,
					Color:     color,
					ShowItems: false,
				})
				chart.Labels = labels
			}
		}
	}

	{
		chart := this.Plugin.Widgets[0].Charts[2].(*teacharts.ProgressBar)
		stat, err := mem.SwapMemory()
		if err != nil {
			logs.Error(err)
		} else {
			chart.Value = stat.UsedPercent
			if chart.Value < 50 {
				chart.Color = teacharts.ColorGreen
			} else if chart.Value > 80 {
				chart.Color = teacharts.ColorRed
			}
			chart.Detail = formatBytes(stat.Used) + "/" + formatBytes(stat.Total)
		}
	}
}
