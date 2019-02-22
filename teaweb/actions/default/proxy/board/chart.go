package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ChartAction actions.Action

// 图表详情
func (this *ChartAction) Run(params struct {
	ServerId string
	WidgetId string
	ChartId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到server")
	}
	this.Data["server"] = maps.Map{
		"id":       server.Id,
		"name":     server.Name,
		"filename": server.Filename,
	}

	widget := widgets.NewWidgetFromId(params.WidgetId)
	if widget == nil {
		this.Fail("找不到Widget")
	}

	this.Data["widget"] = widget
	chart := widget.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	this.Data["chart"] = chart
	this.Show()
}
