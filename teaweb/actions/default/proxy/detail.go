package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type DetailAction actions.Action

func (this *DetailAction) Run(params struct {
	Filename string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if proxy.Index == nil {
		proxy.Index = []string{}
	}

	this.Data["filename"] = params.Filename
	this.Data["proxy"] = proxy

	this.Show()
}
