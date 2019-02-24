package proxy

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method == http.MethodGet && !action.HasPrefix("/proxy/status") {
		proxyutils.AddServerMenu(action)
	}
}
