package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
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
		"id":            waf.Id,
		"name":          waf.Name,
		"on":            waf.On,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
	}

	this.Data["groups"] = lists.Map(teawaf.Template().Inbound, func(k int, v interface{}) interface{} {
		g := v.(*rules.RuleGroup)
		group := waf.FindRuleGroupWithCode(g.Code)

		return maps.Map{
			"name":      g.Name,
			"code":      g.Code,
			"isChecked": group != nil && group.On,
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
	template := teawaf.Template()
	for _, groupCode := range params.GroupCodes {
		g := waf.FindRuleGroupWithCode(groupCode)
		if g != nil {
			g.On = true
			continue
		}
		g = template.FindRuleGroupWithCode(groupCode)
		g.Id = stringutil.Rand(16)
		g.On = true
		waf.AddRuleGroup(g)
	}

	// remove old group {
	for _, g := range waf.Inbound {
		if len(g.Code) > 0 && !lists.ContainsString(params.GroupCodes, g.Code) {
			g.On = false
			continue
		}
	}

	for _, g := range waf.Outbound {
		if len(g.Code) > 0 && !lists.ContainsString(params.GroupCodes, g.Code) {
			g.On = false
			continue
		}
	}

	for _, g := range waf.Inbound {
		logs.Println(g.Code, g.Name, g.Name, g.On)
	}

	filename := "waf." + waf.Id + ".conf"
	err := waf.Save(Tea.ConfigFile(filename))
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(waf.Id) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
