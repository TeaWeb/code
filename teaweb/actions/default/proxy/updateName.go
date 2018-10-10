package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type UpdateNameAction actions.Action

func (this *UpdateNameAction) Run(params struct {
	Filename string
	Index    int
	Name     string
	Must     *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入域名")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Name) {
		proxy.Name[params.Index] = params.Name
	}

	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success("保存成功")
}
