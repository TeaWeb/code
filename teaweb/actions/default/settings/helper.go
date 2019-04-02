package settings

import (
	"github.com/TeaWeb/code/teaweb/configs"
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

	// TabBar菜单
	tabbar := []maps.Map{}

	user := configs.SharedAdminConfig().FindActiveUser(action.Session().GetString("username"))
	if user.Granted(configs.AdminGrantAll) {
		tabbar = append(tabbar, map[string]interface{}{
			"name":    "管理界面",
			"subName": "",
			"url":     "/settings",
			"active":  action.Spec.HasClassPrefix("settings.IndexAction", "server."),
		})
	}

	tabbar = append(tabbar, map[string]interface{}{
		"name":    "个人资料",
		"subName": "",
		"url":     "/settings/profile",
		"active":  action.Spec.HasClassPrefix("profile."),
	})

	tabbar = append(tabbar, map[string]interface{}{
		"name":    "登录设置",
		"subName": "",
		"url":     "/settings/login",
		"active":  action.Spec.HasClassPrefix("login."),
	})

	if user.Granted(configs.AdminGrantAll) {
		// mongodb管理
		tabbar = append(tabbar, map[string]interface{}{
			"name":    "MongoDB",
			"subName": "",
			"url":     "/settings/mongo",
			"active":  action.Spec.HasClassPrefix("mongo."),
		})

		// 备份
		tabbar = append(tabbar, map[string]interface{}{
			"name":    "备份",
			"subName": "",
			"url":     "/settings/backup",
			"active":  action.Spec.HasClassPrefix("backup."),
		})
	}

	tabbar = append(tabbar, map[string]interface{}{
		"name":    "检查版本更新",
		"subName": "",
		"url":     "/settings/update",
		"active":  action.Spec.HasClassPrefix("update."),
	})

	action.Data["teaTabbar"] = tabbar
}
