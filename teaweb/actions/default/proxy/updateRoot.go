package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateRootAction actions.Action

func (this *UpdateRootAction) Run(params struct {
	Filename string
	Root     string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Root = params.Root
	proxy.Save()

	proxyutils.NotifyChange()

	this.Success()
}
