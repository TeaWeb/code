package helpers

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type UserMustAuth struct {
	Username string
}

func (this *UserMustAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	var action = actionPtr.Object()
	var session = action.Session()
	var username = session.GetString("username")
	if len(username) == 0 {
		this.login(action)
		return false
	}

	// 检查用户是否存在
	user := configs.SharedAdminConfig().FindUser(username)
	if user == nil {
		this.login(action)
		return false
	}

	this.Username = username

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	// 初始化变量
	modules := []map[string]interface{}{
		{
			"code":     "proxy",
			"menuName": "代理设置",
			"icon":     "paper plane outline",
		},

		/**{
			"code":     "lab",
			"menuName": "实验室",
			"icon":     "medapps",
		},**/
	}

	if teaconst.PlusEnabled {
		modules = append(modules, []map[string]interface{}{
			{
				"code":     "plus.q",
				"menuName": "测试小Q+",
				"icon":     "dog",
			},
			{
				"code":     "plus.apis",
				"menuName": "API+",
				"icon":     "shekel sign",
			},
			{
				"code":     "plus.team",
				"menuName": "团队+",
				"icon":     "users",
			},
		} ...)
	}

	// 附加功能
	modules = append(modules, []map[string]interface{}{
		{
			"code":     "log",
			"menuName": "日志",
			"icon":     "history",
		},
		{
			"code":     "stat",
			"menuName": "统计",
			"icon":     "chart area",
		},
		{
			"code":     "apps",
			"menuName": "本地服务",
			"icon":     "gem outline",
		},
		{
			"code":     "plugins",
			"menuName": "插件",
			"icon":     "puzzle piece",
		},
	} ...)

	action.Data["teaTitle"] = "TeaWeb管理平台"

	if len(user.Name) == 0 {
		action.Data["teaUsername"] = username
	} else {
		action.Data["teaUsername"] = user.Name
	}

	action.Data["teaUserAvatar"] = user.Avatar

	action.Data["teaMenu"] = ""
	action.Data["teaModules"] = modules
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}
	action.Data["teaVersion"] = teaconst.TeaVersion

	return true
}

func (this *UserMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/login")
}
