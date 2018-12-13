package cache

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"net/http"
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
				"active":  false,
			},
		}
		if lists.Contains([]string{"proxy.IndexAction", "proxy.AddAction", "cache.IndexAction"}, action.Spec.ClassName) {
			tabbar = append(tabbar, maps.Map{
				"name":    "添加新代理",
				"subName": "",
				"url":     "/proxy/add",
				"active":  action.Spec.ClassName == "proxy.AddAction",
			})

			tabbar = append(tabbar, maps.Map{
				"name":    "缓存策略",
				"subName": "",
				"url":     "/cache",
				"active":  action.Spec.HasClassPrefix("cache"),
			})
		}
		action.Data["teaTabbar"] = tabbar
	}
}
