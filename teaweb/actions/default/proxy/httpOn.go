package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type HttpOnAction actions.Action

func (this *HttpOnAction) Run(params struct {
	Filename string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Http = true
	proxy.Save()

	proxyutils.NotifyChange()

	this.Success()
}
