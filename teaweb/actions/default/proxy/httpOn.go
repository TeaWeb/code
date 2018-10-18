package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/TeaWeb/code/teaconfigs"
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
	proxy.WriteBack()

	global.NotifyChange()

	this.Success()
}
