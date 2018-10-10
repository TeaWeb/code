package helpers

import (
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

	this.Username = username

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	// 初始化变量
	action.Data["teaTitle"] = "TeaWeb管理平台"
	action.Data["teaUsername"] = username
	action.Data["teaMenu"] = ""
	action.Data["teaModules"] = []map[string]interface{}{
		{
			"code":     "proxy",
			"menuName": "代理设置",
			"icon":     "compass outline",
		},
		{
			"code":     "log",
			"menuName": "访问日志",
			"icon":     "clock outline",
		},
		{
			"code":     "stat",
			"menuName": "统计",
			"icon":     "chart area",
		},
		/**{
			"code":     "monitor",
			"menuName": "监控",
			"icon":     "bell outline",
		},**/
	}
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}

	return true
}

func (this *UserMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/login")
}
