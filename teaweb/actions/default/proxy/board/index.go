package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/board/scripts"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 面板
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

	configFile := "board." + server.Id + ".conf"
	if files.NewFile(Tea.ConfigFile(configFile)).Exists() {
		this.Data["config"] = configFile
	} else {
		this.Data["config"] = "board.default.conf"
	}

	this.Show()
}

// 面板数据
func (this *IndexAction) RunPost(params struct {
	Server string
	Config string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Server: server,
	})
	err = engine.RunConfig(Tea.ConfigFile(params.Config), maps.Map{})
	this.Data["widgetError"] = ""
	if err != nil {
		this.Data["widgetError"] = err.Error()
	}
	this.Data["charts"] = engine.Charts()
	this.Success()
}
