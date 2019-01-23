package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 设置首页
func (this *IndexAction) Run(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	if agent.IsLocal() {
		this.RedirectURL("/agents/board")
		return
	}

	this.Data["agent"] = agent

	this.Show()
}
