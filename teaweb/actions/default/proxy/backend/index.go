package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 后端列表
func (this *IndexAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "backend"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server

	this.Data["queryParams"] = maps.Map{
		"server": params.Server,
	}

	this.Show()
}
