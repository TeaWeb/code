package settings

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

	action.Data["serverChanged"] = serverChanged

	if action.Spec.HasClassPrefix("profile") {
		action.Data["teaMenu"] = "settings.profile"
	} else {
		action.Data["teaMenu"] = "settings"
	}
	action.Data["teaTabbar"] = []maps.Map{
		{
			"name":    "管理界面",
			"subName": "",
			"url":     "/settings",
			"active":  action.Spec.HasClassPrefix("settings.IndexAction", "server."),
		},
		{
			"name":    "个人资料",
			"subName": "",
			"url":     "/settings/profile",
			"active":  action.Spec.HasClassPrefix("profile."),
		},
		{
			"name":    "登录设置",
			"subName": "",
			"url":     "/settings/login",
			"active":  action.Spec.HasClassPrefix("login."),
		},
		{
			"name":    "MongoDB",
			"subName": "",
			"url":     "/settings/mongo",
			"active":  action.Spec.HasClassPrefix("mongo."),
		},
	}
}
