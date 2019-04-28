package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"

	wafactions "github.com/TeaWeb/code/teawaf/actions"
)

type GroupAction actions.Action

// 分组信息
func (this *GroupAction) RunGet(params struct {
	WafId   string
	GroupId string
	Inbound bool
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
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	this.Data["inbound"] = group.IsInbound
	this.Data["outbound"] = !group.IsInbound

	this.Data["group"] = group

	// rule sets
	this.Data["sets"] = lists.Map(group.RuleSets, func(k int, v interface{}) interface{} {
		set := v.(*rules.RuleSet)
		return maps.Map{
			"id":   set.Id,
			"name": set.Name,
			"rules": lists.Map(set.Rules, func(k int, v interface{}) interface{} {
				rule := v.(*rules.Rule)

				return maps.Map{
					"param":             rule.Param,
					"operator":          rule.Operator,
					"value":             rule.Value,
					"isCaseInsensitive": rule.IsCaseInsensitive,
				}
			}),
			"on":         set.On,
			"action":     strings.ToUpper(set.Action),
			"actionName": wafactions.FindActionName(set.Action),
			"connector":  strings.ToUpper(set.Connector),
		}
	})

	this.Show()
}
