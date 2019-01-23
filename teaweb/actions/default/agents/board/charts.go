package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ChartsAction actions.Action

// 图表列表
func (this *ChartsAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId
	this.Data["tabbar"] = "charts"

	charts := []maps.Map{}

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("无法创建看板的配置文件")
	}

	for _, app := range agent.Apps {
		if !app.On {
			continue
		}

		for _, item := range app.Items {
			if !item.On {
				continue
			}

			if len(item.Charts) == 0 {
				continue
			}

			for _, chart := range item.Charts {
				if !chart.On {
					continue
				}
				charts = append(charts, maps.Map{
					"id":       chart.Id,
					"name":     chart.Name,
					"typeName": widgets.FindChartTypeName(chart.Type),
					"app": maps.Map{
						"id":   app.Id,
						"name": app.Name,
					},
					"item": maps.Map{
						"id":   item.Id,
						"name": item.Name,
					},
					"isUsing": board.FindChart(chart.Id) != nil,
				})
			}
		}
	}

	this.Data["charts"] = charts

	this.Show()
}
