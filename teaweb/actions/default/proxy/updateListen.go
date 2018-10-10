package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type UpdateListenAction actions.Action

func (this *UpdateListenAction) Run(params struct {
	Filename string
	Index    int
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

	if params.Index >= 0 && params.Index < len(proxy.Listen) {
		proxy.Listen[params.Index] = params.Listen
	}

	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success("保存成功")
}
