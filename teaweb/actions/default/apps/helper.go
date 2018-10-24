package apps

import (
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/actions"
	"net/http"
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

	/**tabbar := []maps.Map{
		{
			"name":    "关注",
			"subName": "",
			"url":     "/apps",
			"active":  action.Spec.HasClassPrefix("apps.IndexAction"),
		},
		{
			"name":    "所有",
			"subName": fmt.Sprintf("%d", countApps),
			"url":     "/apps/all",
			"active":  action.Spec.HasClassPrefix("apps.AllAction"),
		},
	}

	action.Data["teaTabbar"] = tabbar**/
}
