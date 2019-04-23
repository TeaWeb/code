package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf/groups"
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
		"id":        waf.Id,
		"name":      waf.Name,
		"countSets": waf.CountRuleSets(),
		"on":        waf.On,
	}

	this.Data["groups"] = lists.Map(groups.InternalGroups, func(k int, v interface{}) interface{} {
		g := v.(*rules.RuleGroup)

		return maps.Map{
			"name":      g.Name,
			"code":      g.Code,
			"isChecked": waf.ContainsGroupCode(g.Code),
		}
	})

	this.Show()
}
