package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type DeleteListenAction actions.Action

func (this *DeleteListenAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Listen) {
		proxy.Listen = lists.Remove(proxy.Listen, params.Index).([]string)
	}

	logs.Println(proxy.Listen)
	proxy.WriteToFilename(params.Filename)

	// 重启服务
	global.NotifyChange()

	this.Refresh().Success()
}
