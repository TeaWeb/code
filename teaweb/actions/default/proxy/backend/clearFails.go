package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type ClearFailsAction actions.Action

// 清除失败次数
func (this *ClearFailsAction) Run(params struct {
	Server     string
	LocationId string
	Websocket  bool
	BackendId  string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	runningServer, _ := teaproxy.FindServer(server.Id)
	if runningServer != nil {
		backendList, _ := runningServer.FindBackendList(params.LocationId, params.Websocket)
		if backendList != nil {
			backend := backendList.FindBackend(params.BackendId)
			if backend != nil {
				backend.CurrentFails = 0
			}
		}
	}

	this.Success()
}
