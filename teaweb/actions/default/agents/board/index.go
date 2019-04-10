package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板首页
func (this *IndexAction) Run(params struct {
	AgentId string
}) {
	if len(params.AgentId) == 0 {
		params.AgentId = "local"
	}

	this.Data["agentId"] = params.AgentId
	this.Data["tabbar"] = "board"

	this.Show()
}

// 数据
func (this *IndexAction) RunPost(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	// 添加默认App
	if agent != nil && !agent.AppsIsInitialized {
		agent.AddDefaultApps()
		err := agent.Save()

		// 通知更新
		if err == nil {
			agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_AGENT", maps.Map{}))
		}
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("无法读取Board配置")
	}

	mongoEnabled := teamongo.Test() == nil
	engine := scripts.NewEngine()
	engine.SetMongo(mongoEnabled)

	if !mongoEnabled {
		this.Data["charts" ] = []interface{}{}
		this.Data["output"] = []string{}
		return
	}

	for _, c := range board.Charts {
		app := agent.FindApp(c.AppId)
		if app == nil || !app.On {
			continue
		}

		item := app.FindItem(c.ItemId)
		if item == nil || !item.On {
			continue
		}

		chart := item.FindChart(c.ChartId)
		if chart == nil || !chart.On {
			continue
		}

		o, err := chart.AsObject()
		if err != nil {
			logs.Error(err)
			continue
		}

		var chartName = chart.Name + "<span class=\"ops\"><a href=\"\" title=\"从看板移除\" onclick=\"return Tea.Vue.removeChart('" + c.AppId + "', '" + c.ItemId + "', '" + c.ChartId + "')\"><i class=\"icon remove small\"></i></a> &nbsp; <a href=\"/agents/apps/itemValues?agentId=" + agent.Id + "&appId=" + app.Id + "&itemId=" + item.Id + "\" title=\"查看数值记录\"><i class=\"icon external small\"></i></a></span>"
		code, err := o.AsJavascript(maps.Map{
			"name":    chartName,
			"columns": chart.Columns,
		})
		if err != nil {
			logs.Error(err)
			continue
		}

		engine.SetContext(&scripts.Context{
			Agent: agent,
			App:   app,
			Item:  item,
		})

		widgetCode := `var widget = new widgets.Widget({
	"name": "看板",
	"requirements": ["mongo"]
});

widget.run = function () {
`
		widgetCode += "{\n" + code + "\n}\n"
		widgetCode += `
};
`

		err = engine.RunCode(widgetCode)
		if err != nil {
			logs.Error(err)
			continue
		}
	}

	this.Data["charts" ] = engine.Charts()
	this.Data["output"] = engine.Output()
	this.Success()
}
