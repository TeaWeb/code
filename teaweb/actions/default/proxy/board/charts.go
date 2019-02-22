package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type ChartsAction actions.Action

// 图表
func (this *ChartsAction) Run(params struct {
	Server string
	Type   string
}) {
	if len(params.Type) == 0 {
		params.Type = "realtime"
	}

	this.Data["boardType"] = params.Type

	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	this.Data["server"] = maps.Map{
		"id":       server.Id,
		"filename": params.Server,
	}

	// 正在使用的图表
	usingCharts := []maps.Map{}
	if params.Type == "realtime" {
		if server.RealtimeBoard != nil {
			for _, c := range server.RealtimeBoard.Charts {
				chart := c.FindChart()
				if chart == nil {
					continue
				}
				usingCharts = append(usingCharts, maps.Map{
					"id":           chart.Id,
					"name":         chart.Name,
					"description":  chart.Description,
					"requirements": chart.Requirements,
					"columns":      chart.Columns,
					"on":           chart.On,
					"widgetId":     c.WidgetId,
				})
			}
		}
	} else {
		if server.StatBoard != nil {
			for _, c := range server.StatBoard.Charts {
				chart := c.FindChart()
				if chart == nil {
					continue
				}
				usingCharts = append(usingCharts, maps.Map{
					"id":           chart.Id,
					"name":         chart.Name,
					"description":  chart.Description,
					"requirements": chart.Requirements,
					"columns":      chart.Columns,
					"on":           chart.On,
					"widgetId":     c.WidgetId,
				})
			}
		}
	}
	this.Data["usingCharts"] = usingCharts

	// 所有的图表
	this.Data["widgets"] = lists.Map(widgets.LoadAllWidgets(), func(k int, v interface{}) interface{} {
		widget := v.(*widgets.Widget)

		return maps.Map{
			"id": widget.Id,
			"charts": lists.Map(widget.Charts, func(k int, v interface{}) interface{} {
				chart := v.(*widgets.Chart)
				isUsing := false
				if params.Type == "realtime" {
					if server.RealtimeBoard != nil {
						isUsing = server.RealtimeBoard.HasChart(chart.Id)
					}
				} else {
					if server.StatBoard != nil {
						isUsing = server.StatBoard.HasChart(chart.Id)
					}
				}
				return maps.Map{
					"id":           chart.Id,
					"name":         chart.Name,
					"description":  chart.Description,
					"requirements": chart.Requirements,
					"columns":      chart.Columns,
					"on":           chart.On,
					"isUsing":      isUsing,
				}
			}),
		}
	})

	this.Show()
}
