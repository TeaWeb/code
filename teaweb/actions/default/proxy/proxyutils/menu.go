package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// 添加服务器菜单
func AddServerMenu(action *actions.ActionObject) {
	// 选中代理
	action.Data["teaMenu"] = "proxy"

	// 子菜单
	var subMenu = utils.NewSubMenu()

	// 服务
	var hasServer = false
	for _, server := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
		urlPrefix := "/proxy/board"
		if action.HasPrefix("/stat") {
			urlPrefix = "/stat"
		} else if action.HasPrefix("/proxy/log") {
			urlPrefix = "/proxy/log"
		} else if action.HasPrefix("/proxy") && !action.HasPrefix("/proxy/board", "/proxy/add") {
			urlPrefix = "/proxy/detail"
		}
		subMenu.Add(server.Description, "", urlPrefix+"?server="+server.Filename, action.ParamString("server") == server.Filename)
		hasServer = true
	}
	if hasServer {
		action.Data["teaSubHeader"] = "代理服务"
	}

	// 其他
	subMenu.Add("[添加新代理]", "", "/proxy/add", action.Spec.ClassName == "proxy.AddAction", )
	subMenu.Add("[缓存策略]", "", "/cache", action.Spec.HasClassPrefix("cache"), )
	utils.SetSubMenu(action, subMenu)

	// Tabbar
	if hasServer {
		if action.Spec.HasClassPrefix(
			"proxy", "ssl", "locations", "fastcgi", "rewrite", "headers", "backend", "board", "websocket",
			"stat.",
			"log.",
		) && !action.Spec.HasClassPrefix("proxy.AddAction", "log.RuntimeAction") {
			serverFilename := action.ParamString("server")
			tabbar := []maps.Map{
				{
					"name":    "看板",
					"subName": "",
					"url":     "/proxy/board?server=" + serverFilename,
					"active":  action.Spec.HasClassPrefix("board."),
					"icon":    "dashboard",
				},
				{
					"name":    "日志",
					"subName": "",
					"url":     "/proxy/log?server=" + serverFilename,
					"active":  action.Spec.HasClassPrefix("log."),
					"icon":    "history",
				},
				{
					"name":    "统计",
					"subName": "",
					"url":     "/stat?server=" + serverFilename,
					"active":  action.Spec.HasClassPrefix("stat."),
					"icon":    "chart area",
				},
				/**{
				"name":    "测试",
				"subName": "",
				"url":     "/test?server=" + serverFilename,
				"active":  false,
				"icon":    "stethoscope",
			},**/
				{
					"name":    "设置",
					"subName": "",
					"url":     "/proxy/detail?server=" + serverFilename,
					"icon":    "setting",
					"active":  action.Spec.HasClassPrefix("proxy", "ssl", "locations", "fastcgi", "rewrite", "headers", "backend", "websocket") && !action.HasPrefix("/proxy/delete"),
				},
				{
					"name":    "删除",
					"subName": "",
					"url":     "/proxy/delete?server=" + serverFilename,
					"icon":    "trash",
					"active":  action.HasPrefix("/proxy/delete"),
				},
			}
			action.Data["teaTabbar"] = tabbar
		}
	}
}
