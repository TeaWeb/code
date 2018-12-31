package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
)

type DeleteNameAction actions.Action

func (this *DeleteNameAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Name) {
		proxy.Name = lists.Delete(proxy.Name, params.Index).([]string)
	}

	proxy.Save()

	// 重启服务
	proxyutils.NotifyChange()

	this.Refresh().Success()
}
