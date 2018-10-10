package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/lists"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method == http.MethodGet {
		action.Data["teaMenu"] = "proxy"

		tabbar := []maps.Map{
			{
				"name":    "已有代理",
				"subName": "",
				"url":     "/proxy",
				"active":  action.Spec.ClassName != "proxy.AddAction",
			},
		}
		if lists.Contains([]string{"proxy.IndexAction", "proxy.AddAction"}, action.Spec.ClassName) {
			tabbar = append(tabbar, maps.Map{
				"name":    "添加新代理",
				"subName": "",
				"url":     "/proxy/add",
				"active":  action.Spec.ClassName == "proxy.AddAction",
			})
		}
		action.Data["teaTabbar"] = tabbar
	}
}
