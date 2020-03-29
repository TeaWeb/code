package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 添加服务器菜单
func AddServerMenu(actionWrapper actions.ActionWrapper) {
	action := actionWrapper.Object()

	// 选中代理
	action.Data["teaMenu"] = "proxy"

	// 子菜单
	menuGroup := utils.NewMenuGroup()

	// 服务
	var hasServer = false
	var isTCP = false

	isIndex := !action.HasPrefix("/proxy/add", "/cache", "/proxy/waf", "/proxy/log/policies", "/proxy/certs", "/proxy/settings")

	serverId := action.ParamString("serverId")
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		logs.Error(err)
	}

	if isIndex {
		menu := menuGroup.FindMenu("", "")
		for _, server := range serverList.FindAllServers() {
			urlPrefix := "/proxy/board"
			if server.IsTCP() { // TCP
				urlPrefix = "/proxy/detail"
			} else { // HTTP
				if action.HasPrefix("/proxy/stat") {
					urlPrefix = "/proxy/stat"
				} else if action.HasPrefix("/proxy/log") {
					urlPrefix = "/proxy/log"
				} else if action.HasPrefix("/proxy") && !action.HasPrefix("/proxy/board", "/proxy/add") {
					urlPrefix = "/proxy/detail"
				}
			}

			item := menu.Add(server.Description, "", urlPrefix+"?serverId="+server.Id, serverId == server.Id)
			item.IsSortable = true

			if server.IsTCP() {
				item.SupName = "tcp"
			} else if server.ForwardHTTP != nil {
				item.SupName = "forward"
			}

			// port
			ports := []string{}
			if server.Http {
				for _, listen := range server.Listen {
					index := strings.LastIndex(listen, ":")
					if index > -1 {
						ports = append(ports, ":"+listen[index+1:])
					} else {
						ports = append(ports, ":"+listen)
					}
				}
			}
			if server.SSL != nil && server.SSL.On {
				for _, listen := range server.SSL.Listen {
					index := strings.LastIndex(listen, ":")
					if index > -1 {
						ports = append(ports, ":"+listen[index+1:])
					} else {
						ports = append(ports, ":"+listen)
					}
				}
			}
			if len(ports) > 0 {
				if len(ports) > 2 {
					item.SubName = ports[0] + ", " + ports[1] + "等 "
				} else {
					item.SubName = strings.Join(ports, ", ") + " "
				}
			}

			// on | off
			if (server.IsHTTP() && !server.Http && (server.SSL == nil || !server.SSL.On)) || (server.IsTCP() && (server.TCP == nil || !server.TCP.TCPOn)) {
				item.SubName = "未启用"
				item.SubColor = "red"
			}

			if server.Id == serverId {
				isTCP = server.IsTCP()
			}

			hasServer = true
		}
		if hasServer {
			if action.Request.URL.Path == "/proxy/board" {
				menu.Name = "代理服务"
			} else {
				menu.Name = "代理服务"
			}
		}
	}

	// 其他
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Add("代理服务", "", "/proxy", isIndex)
		menu.Add("+添加新代理", "", "/proxy/add", action.Spec.ClassName == "proxy.AddAction", )
		menu.Add("缓存策略", "", "/cache", action.HasPrefix("/cache"))
		menu.Add("WAF策略", "", "/proxy/waf", action.HasPrefix("/proxy/waf"))
		menu.Add("日志策略", "", "/proxy/log/policies", action.HasPrefix("/proxy/log/policies"))
		menu.Add("SSL证书管理", "", "/proxy/certs", action.HasPrefix("/proxy/certs"))
		menu.Add("通用设置", "", "/proxy/settings", action.HasPrefix("/proxy/settings"))
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
			"notices.",
		) && !action.Spec.HasClassPrefix("proxy.AddAction", "log.RuntimeAction") {
			if isTCP { // TCP
				tabbar := []maps.Map{
					{
						"name":    "设置",
						"subName": "",
						"url":     "/proxy/detail?serverId=" + serverId,
						"icon":    "setting",
						"active":  !action.HasPrefix("/proxy/delete"),
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
			} else { // HTTP
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
						"active":  action.Spec.HasClassPrefix("proxy", "ssl", "locations", "fastcgi", "rewrite", "headers", "backend", "websocket", "access", "servers", "tunnel", "notices") && !action.HasPrefix("/proxy/delete"),
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
