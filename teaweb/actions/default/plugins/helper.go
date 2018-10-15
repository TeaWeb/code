package plugins

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "plugins"
	action.Data["teaTabbar"] = []maps.Map{
		{
			"name":    "已安装插件",
			"subName": "",
			"url": "/plugins",
			"active":  true,
		},
	}
}
