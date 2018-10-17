package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type DeleteIgnoreAction actions.Action

func (this *DeleteIgnoreAction) Run(params struct {
	Filename string
	Index    int
	Must     *actions.Must
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.DeleteIgnoreHeaderAtIndex(params.Index)
	proxy.WriteBack()

	global.NotifyChange()

	this.Success()
}
