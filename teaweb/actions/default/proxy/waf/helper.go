package waf

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

// 缓存相关Helper
func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method == http.MethodGet {
		proxyutils.AddServerMenu(action)

		action.Data["selectedMenu"] = "list"
		if action.HasPrefix("/proxy/waf/add") {
			action.Data["selectedMenu"] = "add"
		}

		action.Data["selectedSubMenu"] = "detail"
		if action.HasPrefix("/proxy/waf/rules","/proxy/waf/group") {
			action.Data["selectedSubMenu"] = "rules"
		}
	}
}
