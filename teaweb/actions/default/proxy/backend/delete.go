package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除后端服务器
func (this *DeleteAction) Run(params struct {
	Server     string
	LocationId string
	Websocket  bool
	BackendId  string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}
	backendList.DeleteBackend(params.BackendId)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Refresh().Success()
}
