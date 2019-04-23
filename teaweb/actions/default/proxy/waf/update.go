package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf/groups"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) RunGet(params struct {
	WafId string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = maps.Map{
		"id":        waf.Id,
		"name":      waf.Name,
		"On":        waf.On,
		"countSets": waf.CountRuleSets(),
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

// 保存修改
func (this *UpdateAction) RunPost(params struct {
	WafId      string
	Name       string
	GroupCodes []string
	On         bool
	Must       *actions.Must
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("无法找到WAF")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	waf.Name = params.Name
	waf.On = params.On

	// add new group
	for _, groupCode := range params.GroupCodes {
		if waf.ContainsGroupCode(groupCode) {
			continue
		}
		for _, g := range groups.InternalGroups {
			if g.Code == groupCode {
				newGroup := rules.NewRuleGroup()
				newGroup.Id = stringutil.Rand(16)
				newGroup.On = g.On
				newGroup.Code = g.Code
				newGroup.Name = g.Name
				newGroup.RuleSets = g.RuleSets
				waf.AddRuleGroup(newGroup)
			}
		}
	}

	// remove old group
	result := []*rules.RuleGroup{}
	for _, g := range waf.RuleGroups {
		if len(g.Code) > 0 && !lists.ContainsString(params.GroupCodes, g.Code) {
			continue
		}
		result = append(result, g)
	}
	waf.RuleGroups = result

	filename := "waf." + waf.Id + ".conf"
	err := waf.Save(Tea.ConfigFile(filename))
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
