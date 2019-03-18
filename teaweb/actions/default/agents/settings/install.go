package settings

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type InstallAction actions.Action

// 安装部署
func (this *InstallAction) Run(params struct {
	AgentId string
}) {
	this.Data["selectedTab"] = "install"

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
