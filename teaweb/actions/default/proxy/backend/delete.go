package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除后端服务器
func (this *DeleteAction) Run(params struct {
	Server    string
	BackendId string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	server.DeleteBackend(params.BackendId)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	global.NotifyChange()

	this.Refresh().Success()
}
