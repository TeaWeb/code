package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板
func (this *IndexAction) Run(params struct {
	ServerId  string
	BoardType string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要查看的代理服务")
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	this.Data["boardType"] = params.BoardType

	this.Show()
}

// 看板数据
func (this *IndexAction) RunPost(params struct {
	ServerId string
	Type     string // realtime or stat
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要查看的代理服务")
	}

	if len(params.Type) == 0 {
		params.Type = "realtime"
	}

	var board *teaconfigs.Board
	switch params.Type {
	case "realtime":
		board = server.RealtimeBoard
	case "stat":
		board = server.StatBoard
	}

	// 初始化
	if board == nil {
		board = teaconfigs.NewBoard()
		switch params.Type {
		case "realtime":
			server.RealtimeBoard = board

			// 添加一些默认的图表
			board.AddChart("teaweb.proxy_status", "kTVuOEBm605H3AJS")
			board.AddChart("teaweb.locations", "cyZsvwR66oxcVpvj")
			board.AddChart("teaweb.bandwidth_realtime", "hiAUsteL1V6LG8zD")
			board.AddChart("teaweb.request_realtime", "APvSaVEoQ7VUvX4a")
			board.AddChart("teaweb.request_time", "g8SrxuMYwwhxNWwk")
			board.AddChart("teaweb.status_stat", "xnUsgQMSjWZ9MN7g")
			board.AddChart("teaweb.latest_errors", "RUCF1EbF4FpPMHpN")
			err := server.Save()
			if err != nil {
				logs.Error(err)
			}
		case "stat":
			server.StatBoard = board

			// TODO 添加一些默认的图表
		}
	}

	if len(board.Charts) == 0 {
		this.Data["charts"] = []maps.Map{}
		this.Success()
	}

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Server: server,
	})

	for _, c := range board.Charts {
		_, chart := c.FindChart()
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
