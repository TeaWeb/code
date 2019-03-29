package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateItemChartAction actions.Action

// 给监控项添加图标
func (this *UpdateItemChartAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	ChartId string
	From    string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	chart := item.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	this.Data["from"] = params.From
	this.Data["item"] = item
	this.Data["chart"] = chart
	this.Data["chartTypes"] = widgets.AllChartTypes

	source := agents.FindDataSource(item.SourceCode)
	if source != nil {
		this.Data["selectedSource"] = maps.Map{
			"variables": source["instance"].(agents.SourceInterface).Variables(),
		}
	} else {
		this.Data["selectedSource"] = nil
	}

	this.Show()
}

// 提交保存
func (this *UpdateItemChartAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string
	ChartId string

	Name      string
	On        bool
	Columns   uint8
	ChartType string

	HTMLCode string `alias:"htmlCode"`

	PieParam string
	PieLimit int

	LineParams []string
	LineLimit  int

	URL string `alias:"urlURL"`

	JavascriptCode string

	Must *actions.Must
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到要修改的App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到要操作的Item")
	}

	chart := item.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	chart.Name = params.Name
	chart.On = params.On
	chart.Columns = params.Columns
	chart.Type = params.ChartType

	switch params.ChartType {
	case "html":
		options := &widgets.HTMLChart{}
		options.HTML = params.HTMLCode
		teautils.ObjectToMapJSON(options, &chart.Options)
	case "url":
		options := &widgets.URLChart{}
		options.URL = params.URL
		teautils.ObjectToMapJSON(options, &chart.Options)
	case "pie":
		options := &widgets.PieChart{}
		options.Param = params.PieParam
		options.Limit = params.PieLimit
		teautils.ObjectToMapJSON(options, &chart.Options)
	case "line":
		options := &widgets.LineChart{}
		options.Params = params.LineParams
		options.Limit = params.LineLimit
		teautils.ObjectToMapJSON(options, &chart.Options)
	case "javascript":
		options := &widgets.JavascriptChart{}
		options.Code = params.JavascriptCode

		// 测试
		engine := scripts.NewEngine()
		engine.SetMongo(teamongo.Test() == nil)
		engine.SetContext(&scripts.Context{
			Agent: agent,
			App:   app,
			Item:  item,
		})
		widgetCode := `var widget = new widgets.Widget({
	
});

widget.run = function () {
` + options.Code + `
};
`
		err := engine.RunCode(widgetCode)
		if err != nil {
			this.Fail("Javascript代码错误：" + err.Error())
		}
		if len(engine.Charts()) == 0 {
			this.Fail("代码中应该包含至少一个图表")
		}

		teautils.ObjectToMapJSON(options, &chart.Options)
	}

	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, nil, nil)
	}

	this.Success()
}
