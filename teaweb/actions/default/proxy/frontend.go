package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type FrontendAction actions.Action

func (this *FrontendAction) Run(params struct {
	Filename string
}) {
	this.Data["selectedTab"] = "frontend"
	this.Data["filename"] = params.Filename

	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["proxy"] = server

	this.Show()
}
