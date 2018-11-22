package probes

import (
	"fmt"
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/logs"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"time"
)

type CPUProbe struct {
	Probe

	cpuValues []float64
}

func (this *CPUProbe) Run() {
	this.InitOnce(func() {
		widget := teaplugins.NewWidget()
		widget.Group = teaplugins.WidgetGroupSystem
		widget.Dashboard = true
		widget.Name = "CPU使用情况"

		{
			chart := teacharts.NewLineChart()
			chart.Max = 100
			chart.YTickCount = 1
			widget.AddChart(chart)
		}

		{
			chart := teacharts.NewTable()
			chart.Name = "负载"
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

	stat, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		logs.Error(err)
	} else if len(stat) > 0 {
		usage := stat[0]
		this.cpuValues = append(this.cpuValues, usage)
		if len(this.cpuValues) > 20 {
			this.cpuValues = this.cpuValues[len(this.cpuValues)-20:]
		}

		values := []interface{}{}
		labels := []string{}
		for _, value := range this.cpuValues {
			values = append(values, value)
			labels = append(labels, "")
		}

		color := teacharts.ColorBlue
		if usage < 50 {
			color = teacharts.ColorGreen
		} else if usage > 80 {
			color = teacharts.ColorRed
		}

		chart := this.Plugin.Widgets[0].Charts[0].(*teacharts.LineChart)
		chart.ResetLines()
		chart.Labels = labels
		chart.AddLine(&teacharts.Line{
			Name:   "实时CPU使用",
			Values: values,
			Filled: true,
			Color:  color,
		})
	}

	{
		stat, err := load.Avg()
		if err != nil {
			logs.Error(err)
		} else {
			chart := this.Plugin.Widgets[0].Charts[1].(*teacharts.Table)
			chart.ResetRows()
			chart.AddRow("负载：", fmt.Sprintf("%.2f", stat.Load1), fmt.Sprintf("%.2f", stat.Load5), fmt.Sprintf("%.2f", stat.Load15))
		}
	}
}
