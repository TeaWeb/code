package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type MoveUpAction actions.Action

func (this *MoveUpAction) Run(params struct {
	ServerId string
	Index    int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if params.Index >= 1 && params.Index < len(server.Locations) {
		prev := server.Locations[params.Index-1]
		current := server.Locations[params.Index]
		server.Locations[params.Index-1] = current
		server.Locations[params.Index] = prev
	}

	err := server.Save()
	if err != nil {
		this.Fail("找不到Server")
	}

	proxyutils.NotifyChange()

	this.Refresh().Success()
}
