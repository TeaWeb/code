package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type AddIgnoreAction actions.Action

func (this *AddIgnoreAction) Run(params struct {
	Filename string
	Name     string
	Must     *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入Name")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.AddIgnoreHeader(params.Name)
	proxy.Save()

	global.NotifyChange()

	this.Refresh().Success()
}
