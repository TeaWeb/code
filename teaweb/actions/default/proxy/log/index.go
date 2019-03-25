package log

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	ServerId string
	LogType  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err := teamongo.Test()
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接MongoDB"
	}

	this.Data["server"] = server

	this.Data["logType"] = params.LogType

	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	proxyutils.AddServerMenu(this)

	this.Show()
}
