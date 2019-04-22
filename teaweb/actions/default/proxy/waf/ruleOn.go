package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type RuleOnAction actions.Action

// 启用规则集
func (this *RuleOnAction) RunPost(params struct {
	WafId   string
	GroupId string
	SetId   string
}) {
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	set := group.FindRuleSet(params.SetId)
	if set == nil {
		this.Fail("找不到规则集")
	}
	set.On = true

	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
