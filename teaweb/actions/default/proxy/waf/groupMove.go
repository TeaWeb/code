package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type GroupMoveAction actions.Action

// 移动分组
func (this *GroupMoveAction) RunPost(params struct {
	WafId     string
	FromIndex int
	ToIndex   int
}) {
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	waf.MoveRuleGroup(params.FromIndex, params.ToIndex)
	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
