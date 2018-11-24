package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
	"regexp"
)

type UpdateIndexAction actions.Action

func (this *UpdateIndexAction) Run(params struct {
	Filename string
	Index    string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if len(params.Index) > 0 {
		indexes := regexp.MustCompile("\\s+").Split(params.Index, -1)
		proxy.Index = indexes
	} else {
		proxy.Index = []string{}
	}
	err = proxy.Save()
	if err != nil {
		this.Fail(err.Error())
	}

	global.NotifyChange()

	this.Data["indexes"] = proxy.Index

	this.Success()
}
