package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type UpdateIgnoreAction actions.Action

func (this *UpdateIgnoreAction) Run(params struct {
	Filename string
	Index    int
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

	proxy.UpdateIgnoreHeaderAtIndex(params.Index, params.Name)
	proxy.WriteBack()

	global.NotifyChange()

	this.Success()
}
