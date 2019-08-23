package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
)

type PreviewItemChartAction actions.Action

// 预览图表
func (this *PreviewItemChartAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string

	Name      string
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

	chart := widgets.NewChart()
	chart.Name = params.Name
	chart.On = true
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
		teautils.ObjectToMapJSON(options, &chart.Options)
	}

	c, err := chart.AsObject()
	if err != nil {
		this.Fail("发现错误：" + err.Error())
	}

	code, err := c.AsJavascript(map[string]interface{}{
		"name":    params.Name,
		"columns": params.Columns,
	})

	mongoEnabled := teamongo.Test() == nil
	engine := scripts.NewEngine()
	engine.SetMongo(mongoEnabled)
	engine.SetCache(false)

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
		this.Fail("发生错误：" + err.Error())
	}

	this.Data["charts" ] = engine.Charts()
	this.Data["output"] = engine.Output()
	this.Success()
}
