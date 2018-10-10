package apps

import (
	"github.com/TeaWeb/code/teaservices/probes"
	"github.com/iwind/TeaGo/logs"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/TeaWeb/code/teacharts"
	"fmt"
)

type MySQLProbe struct {
	probes.Probe
}

func (this *MySQLProbe) Run() {
	this.InitOnce(func() {
		logs.Println("probe mysql")

		widget := teaplugins.NewWidget()
		widget.Dashboard = true
		widget.Group = teaplugins.WidgetGroupService
		widget.Name = "MySQL"
		widget.OnForceReload(func() {
			this.Run()
		})
		this.Plugin.AddWidget(widget)
	})

	widget := this.Plugin.Widgets[0]

	result := ps("mysql", []string{"mysqld_safe"}, false)
	widget.ResetCharts()
	if len(result) == 0 {
		return
	}
	for _, proc := range result {
		chart := teacharts.NewTable()
		chart.AddRow("PID:", fmt.Sprintf("%d", proc.Pid), "<i class=\"ui icon dot circle green\"></i>")
		chart.SetWidth(25, 60, 15)
		widget.AddChart(chart)
	}
}
