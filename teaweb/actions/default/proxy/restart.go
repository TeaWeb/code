package proxy

import (
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type RestartAction actions.Action

func (this *RestartAction) Run(params struct{}) {
	teaproxy.Restart()

	proxyutils.FinishChange()

	this.Success()
}
