package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
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

	proxy.WriteToFilename(params.Filename)

	// 重启服务
	global.NotifyChange()

	this.Refresh().Success()
}
