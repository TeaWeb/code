package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type AddListenAction actions.Action

// 添加监听地址
func (this *AddListenAction) Run(params struct {
	Filename string
	Listen   string
	Must     *actions.Must
}) {
	params.Must.
		Field("listen", params.Listen).
		Require("请输入监听地址")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.AddListen(params.Listen)
	proxy.Save()

	proxyutils.NotifyChange()

	this.Refresh().Success("保存成功")
}
