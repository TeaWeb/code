package proxy

import (
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type RestartAction actions.Action

func (this *RestartAction) Run(params struct{}) {
	err := teaproxy.SharedManager.Restart()
	if err != nil {
		this.Fail("重启失败：" + err.Error())
	}

	proxyutils.FinishChange()

	this.Refresh().Success()
}
