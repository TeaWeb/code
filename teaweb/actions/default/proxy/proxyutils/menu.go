package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

// 添加服务器菜单
func AddServerMenu(actionWrapper actions.ActionWrapper) {
	action := actionWrapper.Object()

	// 选中代理
	action.Data["teaMenu"] = "proxy"

	// 子菜单
	var subMenu = utils.NewSubMenu()

	// 服务
	var hasServer = false
	serverId := action.ParamString("serverId")
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		logs.Error(err)
	}
	for _, server := range serverList.FindAllServers() {
		urlPrefix := "/proxy/board"
		if action.HasPrefix("/proxy/stat") {
			urlPrefix = "/proxy/stat"
		} else if action.HasPrefix("/proxy/log") {
			urlPrefix = "/proxy/log"
		} else if action.HasPrefix("/proxy") && !action.HasPrefix("/proxy/board", "/proxy/add") {
			urlPrefix = "/proxy/detail"
		}
		m := subMenu.Add(server.Description, "", urlPrefix+"?serverId="+server.Id, serverId == server.Id)
		m["sortable"] = true

		hasServer = true
	}
	if hasServer {
		if action.Request.URL.Path == "/proxy/board" {
			action.Data["teaSubHeader"] = "代理服务 <span>(可拖动排序)</span>"
		} else {
			action.Data["teaSubHeader"] = "代理服务"
		}
	}

	// 其他
	subMenu.Add("[添加新代理]", "", "/proxy/add", action.Spec.ClassName == "proxy.AddAction", )
	subMenu.Add("[缓存策略]", "", "/cache", action.Spec.HasClassPrefix("cache"), )
	utils.SetSubMenu(action, subMenu)

	// Tabbar
	if hasServer {
		if action.Spec.HasClassPrefix(
			"proxy",
			"ssl",
			"locations",
			"access.",
			"fastcgi",
			"rewrite",
			"headers",
			"backend",
			"board",
			"websocket",
			"stat.",
			"log.",
		) && !action.Spec.HasClassPrefix("proxy.AddAction", "log.RuntimeAction") {
			tabbar := []maps.Map{
				{
					"name":    "看板",
					"subName": "",
					"url":     "/proxy/board?serverId=" + serverId,
					"active":  action.HasPrefix("/proxy/board") && action.ParamString("boardType") != "stat",
					"icon":    "dashboard",
				},
				{
					"name":    "日志",
					"subName": "",
					"url":     "/proxy/log?serverId=" + serverId,
					"active":  action.HasPrefix("/proxy/log"),
					"icon":    "history",
				},
				{
					"name":    "统计",
					"subName": "",
					"url":     "/proxy/stat?serverId=" + serverId,
					"active":  action.HasPrefix("/proxy/stat") || (action.HasPrefix("/proxy/board") && action.ParamString("boardType") == "stat"),
					"icon":    "chart area",
				},
				{
					"name":    "设置",
					"subName": "",
					"url":     "/proxy/detail?serverId=" + serverId,
					"icon":    "setting",
					"active":  action.Spec.HasClassPrefix("proxy", "ssl", "locations", "fastcgi", "rewrite", "headers", "backend", "websocket", "access") && !action.HasPrefix("/proxy/delete"),
				},
				{
					"name":    "删除",
					"subName": "",
					"url":     "/proxy/delete?serverId=" + serverId,
					"icon":    "trash",
					"active":  action.HasPrefix("/proxy/delete"),
				},
			}
			action.Data["teaTabbar"] = tabbar
		}
	}
}
