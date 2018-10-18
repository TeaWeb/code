package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Locations) {
		proxy.Locations = lists.Remove(proxy.Locations, params.Index).([]*teaconfigs.LocationConfig)
	}

	proxy.WriteBack()

	global.NotifyChange()

	this.Refresh().Success()
}
