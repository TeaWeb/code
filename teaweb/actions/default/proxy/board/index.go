package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板
func (this *IndexAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	this.Data["server"] = maps.Map{
		"id":       server.Id,
		"filename": params.Server,
	}

	this.Show()
}

// 看板数据
func (this *IndexAction) RunPost(params struct {
	Server string
	Type   string // realtime or stat
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	var board *teaconfigs.Board
	{
	}
	switch params.Type {
	case "realtime":
		board = server.RealtimeBoard
	case "stat":
		board = server.StatBoard
	default:
		board = server.RealtimeBoard
	}

	if board == nil || len(board.Charts) == 0 {
		this.Data["charts"] = []maps.Map{}
		this.Success()
	}

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Server: server,
	})

	for _, c := range board.Charts {
		chart := c.FindChart()
		if chart == nil || !chart.On {
			continue
		}

		obj, err := chart.AsObject()
		if err != nil {
			this.Fail(err.Error())
		}
		code, err := obj.AsJavascript(map[string]interface{}{
			"name":    chart.Name,
			"columns": chart.Columns,
		})
		if err != nil {
			this.Fail(err.Error())
		}

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
	}

	this.Data["charts"] = engine.Charts()
	this.Data["output"] = engine.Output()

	this.Success()
}
