package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type DetailAction actions.Action

// 代理详情
func (this *DetailAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	if server.Index == nil {
		server.Index = []string{}
	}

	this.Data["selectedTab"] = "basic"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server

	this.Show()
}
