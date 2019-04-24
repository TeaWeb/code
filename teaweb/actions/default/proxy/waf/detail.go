package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf/inbound"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 详情
func (this *DetailAction) RunGet(params struct {
	WafId string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
		"on":            waf.On,
	}

	this.Data["groups"] = lists.Map(inbound.InternalGroups, func(k int, v interface{}) interface{} {
		g := v.(*rules.RuleGroup)

		return maps.Map{
			"name":      g.Name,
			"code":      g.Code,
			"isChecked": waf.ContainsGroupCode(g.Code),
		}
	})

	// 正在使用此缓存策略的项目
	configItems := []maps.Map{}
	serverList, _ := teaconfigs.SharedServerList()
	if serverList != nil {
		for _, server := range serverList.FindAllServers() {

			if server.WafId == waf.Id {
				configItems = append(configItems, maps.Map{
					"type":   "server",
					"server": server.Description,
					"link":   "/proxy/servers/waf?serverId=" + server.Id,
				})
			}

			for _, location := range server.Locations {
				if location.WafId == waf.Id {
					configItems = append(configItems, maps.Map{
						"type":     "location",
						"server":   server.Description,
						"location": location.Pattern,
						"link":     "/proxy/locations/waf?serverId=" + server.Id + "&locationId=" + location.Id,
					})
				}
			}
		}
	}

	this.Data["configItems"] = configItems

	this.Show()
}
