package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type SetAction actions.Action

func (this *SetAction) Run(params struct {
	Filename string
	Name     string
	Value    string
	Must     *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入Name")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.SetHeader(params.Name, params.Value)
	proxy.Save()

	global.NotifyChange()

	this.Refresh().Success()
}
