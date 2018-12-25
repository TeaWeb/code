package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type FrontendAction actions.Action

// 前端设置
func (this *FrontendAction) Run(params struct {
	Server string
}) {
	this.Data["selectedTab"] = "frontend"
	this.Data["filename"] = params.Server

	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["proxy"] = server

	this.Show()
}
