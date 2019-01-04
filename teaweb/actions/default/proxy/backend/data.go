package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/scheduling"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type DataAction actions.Action

// 后端服务器数据
func (this *DataAction) Run(params struct {
	Server     string
	LocationId string
	Websocket  bool
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}
	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	normalBackends := []*teaconfigs.BackendConfig{}
	backupBackends := []*teaconfigs.BackendConfig{}
	runningServer, _ := teaproxy.FindServer(server.Id)
	for _, backend := range backendList.AllBackends() {
		// 是否下线
		if runningServer != nil {
			runningBackend := runningServer.FindBackend(backend.Id)
			if runningBackend != nil {
				backend.IsDown = runningBackend.IsDown
				backend.DownTime = runningBackend.DownTime
				backend.CurrentFails = runningBackend.CurrentFails
			}
		}

		if backend.IsBackup {
			backupBackends = append(backupBackends, backend)
		} else {
			normalBackends = append(normalBackends, backend)
		}
	}

	this.Data["normalBackends"] = normalBackends
	this.Data["backupBackends"] = backupBackends

	// 算法
	schedulingConfig := backendList.SchedulingConfig()
	if schedulingConfig == nil {
		this.Data["scheduling"] = scheduling.FindSchedulingType("random")
	} else {
		s := scheduling.FindSchedulingType(schedulingConfig.Code)
		if s == nil {
			this.Data["scheduling"] = scheduling.FindSchedulingType("random")
		} else {
			this.Data["scheduling"] = s
		}
	}

	this.Success()
}
