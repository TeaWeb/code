package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/actions"
)

type DetailAction actions.Action

func (this *DetailAction) Run(params struct {
	Server string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	if proxy.Index == nil {
		proxy.Index = []string{}
	}

	this.Data["selectedTab"] = "basic"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = proxy

	// 字符集
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}
