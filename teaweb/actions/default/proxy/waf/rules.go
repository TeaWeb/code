package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type RulesAction actions.Action

// 规则
func (this *RulesAction) RunGet(params struct {
	WafId string
}) {
	config := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if config == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = config
	this.Data["groups"] = lists.Map(config.RuleGroups, func(k int, v interface{}) interface{} {
		group := v.(*rules.RuleGroup)
		return maps.Map{
			"id":            group.Id,
			"code":          group.Code,
			"name":          group.Name,
			"on":            group.On,
			"countRuleSets": len(group.RuleSets),
		}
	})

	this.Show()
}
