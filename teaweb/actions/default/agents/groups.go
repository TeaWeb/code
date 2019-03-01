package agents

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type GroupsAction actions.Action

// 分组管理
func (this *GroupsAction) Run(params struct{}) {
	this.Data["groups"] = agents.SharedGroupConfig().Groups

	this.Show()
}
