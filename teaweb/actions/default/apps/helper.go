package apps

import (
	"fmt"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"runtime"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "apps"

	countApps := 0
	for _, plugin := range teaplugins.Plugins() {
		countApps += len(plugin.Apps)
	}

	if runtime.GOOS == "windows" {
		tabbar := []maps.Map{
			{
				"name":    "发现",
				"subName": fmt.Sprintf("%d", countApps),
				"url":     "/apps/all",
				"active":  action.Spec.HasClassPrefix("apps.AllAction"),
			},
		}

		action.Data["teaTabbar"] = tabbar
	} else {
		tabbar := []maps.Map{
			{
				"name":    "发现",
				"subName": fmt.Sprintf("%d", countApps),
				"url":     "/apps/all",
				"active":  action.Spec.HasClassPrefix("apps.AllAction"),
			},
			{
				"name":    "探针",
				"subName": "",
				"url":     "/apps/probes",
				"active":  action.Spec.HasClassPrefix("apps.ProbesAction", "apps.AddProbeAction", "apps.ProbeScriptAction"),
			},
		}

		action.Data["teaTabbar"] = tabbar
	}
}
