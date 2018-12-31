package proxy

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type StatusAction actions.Action

func (this *StatusAction) Run() {
	this.Data["changed"] = proxyutils.ProxyIsChanged()
	this.Success()
}
