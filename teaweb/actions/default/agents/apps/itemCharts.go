package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type ItemChartsAction actions.Action

// 图表
func (this *ItemChartsAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")
	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	this.Data["item"] = item
	this.Data["intervalSeconds"] = item.IntervalDuration().Seconds()

	source := item.Source()
	if source != nil {
		this.Data["sourceName"] = source.Name()
		this.Data["hasDefaultCharts"] = len(source.Charts()) > 0
	} else {
		this.Data["sourceName"] = ""
		this.Data["hasDefaultCharts"] = false
	}

	this.Show()
}

// 获取图表数据
func (this *ItemChartsAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string
	From    string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agentutils.FindAgentApp(agent, params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	widgetCode := `var widget = new widgets.Widget({
	
});

widget.run = function () {
`

	for _, c := range item.Charts {
		o, err := c.AsObject()
		if err != nil {
			logs.Error(err)
			continue
		}

		var options map[string]interface{}
		if app.IsSystem {
			options = map[string]interface{}{
				"name":    c.Name,
				"columns": c.Columns,
			}
		} else {
			options = map[string]interface{}{
				"name":    c.Name + "<span class=\"ops\"><a href=\"/agents/apps/updateItemChart?agentId=" + params.AgentId + "&appId=" + params.AppId + "&itemId=" + params.ItemId + "&chartId=" + c.Id + "&from=" + params.From + "\" title=\"修改\"><i class=\"icon pencil\"></i></a> &nbsp;<a href=\"\" onclick=\"return Tea.Vue.deleteChart('" + c.Id + "')\" title=\"删除\"><i class=\"icon remove\"></i></a></span>",
				"columns": c.Columns,
			}
		}
		code, err := o.AsJavascript(options)

		if err != nil {
			logs.Error(err)
			continue
		}

		widgetCode += "{\n" + code + "\n}\n"
	}

	widgetCode += `
};
`

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Agent: agent,
		App:   app,
		Item:  item,
	})
	err := engine.RunCode(widgetCode)
	if err != nil {
		logs.Error(err)
	}

	this.Data["charts"] = engine.Charts()
	this.Data["output"] = engine.Output()
	this.Success()
}
