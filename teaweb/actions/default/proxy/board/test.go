package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type TestAction actions.Action

// 测试
func (this *TestAction) Run(params struct {
	ServerId       string
	Name           string
	Description    string
	Columns        uint8
	Items          []string
	JavascriptCode string
	On             bool
	Must           *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	chart := widgets.NewChart()
	chart.On = params.On
	chart.Name = params.Name
	chart.Description = params.Description
	chart.Columns = params.Columns
	chart.Requirements = params.Items
	chart.Type = "javascript"
	chart.Options = maps.Map{
		"code": params.JavascriptCode,
	}
	obj, err := chart.AsObject()
	if err != nil {
		this.Fail("运行错误：" + err.Error())
	}

	code, err := obj.AsJavascript(map[string]interface{}{
		"name":    params.Name,
		"columns": params.Columns,
	})
	if err != nil {
		this.Fail("运行错误：" + err.Error())
	}

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Server: server,
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
		this.Fail("运行错误：" + err.Error())
	}

	this.Data["charts"] = engine.Charts()
	this.Data["output"] = engine.Output()

	this.Success()
}
