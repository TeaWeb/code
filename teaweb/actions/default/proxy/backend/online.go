package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type OnlineAction actions.Action

// 上线服务器
func (this *OnlineAction) Run(params struct {
	Server    string
	BackendId string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	runningServer, _ := teaproxy.FindServer(server.Id)
	if runningServer != nil {
		backend := runningServer.FindBackend(params.BackendId)
		if backend != nil {
			backend.IsDown = false
			backend.CurrentFails = 0
			runningServer.SetupScheduling(false)
		}
	}

	this.Success()
}
