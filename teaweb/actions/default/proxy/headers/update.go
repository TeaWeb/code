package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

func (this *UpdateAction) Run(params struct {
	Filename string
	Index    int
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

	header := proxy.HeaderAtIndex(params.Index)
	if header != nil {
		header.Name = params.Name
		header.Value = params.Value
		proxy.Save()

		global.NotifyChange()
	}

	this.Refresh().Success()
}
