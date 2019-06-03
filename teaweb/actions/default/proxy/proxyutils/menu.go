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
	menuGroup := utils.NewMenuGroup()
	menu := menuGroup.FindMenu("", "")

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
		item := menu.Add(server.Description, "", urlPrefix+"?serverId="+server.Id, serverId == server.Id)
		item.IsSortable = true

		hasServer = true
	}
	if hasServer {
		if action.Request.URL.Path == "/proxy/board" {
			menu.Name = "代理服务"
		} else {
			menu.Name = "代理服务"
		}
	}

	// 其他
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menu.Add("[添加新代理]", "", "/proxy/add", action.Spec.ClassName == "proxy.AddAction", )
		menu.Add("[缓存策略]", "", "/cache", action.HasPrefix("/cache"))
		item := menu.Add("[WAF策略]", "", "/proxy/waf", action.HasPrefix("/proxy/waf"))
		item.SupName = "beta"
	}
	utils.SetSubMenu(action, menuGroup)

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
			"groups",
			"stat.",
			"log.",
			"servers.",
			"tunnel.",
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
					"active":  action.Spec.HasClassPrefix("proxy", "ssl", "locations", "fastcgi", "rewrite", "headers", "backend", "websocket", "access", "servers", "tunnel") && !action.HasPrefix("/proxy/delete"),
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

// 包装Server相关数据
func WrapServerData(server *teaconfigs.ServerConfig) maps.Map {
	return maps.Map{
		"id":          server.Id,
		"description": server.Description,
		"backends":    server.Backends,
		"headers":     server.Headers,
		"cachePolicy": server.CachePolicy,
		"cacheOn":     server.CacheOn,
		"ssl":         server.SSL,
		"locations":   server.Locations,
		"wafOn":       server.WAFOn,
		"wafId":       server.WafId,
		"tunnelOn":    server.Tunnel != nil && server.Tunnel.On,
	}
}
