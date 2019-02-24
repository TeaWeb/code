package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) Run(params struct {
	ServerId string
}) {
	this.Data["server"] = maps.Map{
		"id": params.ServerId,
	}

	this.Show()
}

func (this *DeleteAction) RunPost(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	err := server.Delete()
	if err != nil {
		logs.Error(err)
		this.Fail("配置文件删除失败")
	}

	// @TODO 删除对应的certificate file和certificate key file

	// 重启
	proxyutils.NotifyChange()

	this.Success()
}
