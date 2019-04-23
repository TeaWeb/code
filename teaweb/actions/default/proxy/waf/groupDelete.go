package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/actions"
)

type GroupDeleteAction actions.Action

// 删除分组
func (this *GroupDeleteAction) RunPost(params struct {
	WafId   string
	GroupId string
}) {
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	waf.RemoveRuleGroup(params.GroupId)
	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(waf.Id) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
