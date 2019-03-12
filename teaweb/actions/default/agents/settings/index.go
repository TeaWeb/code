package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
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

	state, isWaiting := agentutils.CheckAgentIsWaiting(agent.Id)
	if isWaiting {
		this.Data["agentVersion"] = state.Version
		this.Data["agentSpeed"] = state.Speed
		this.Data["agentIsWaiting"] = true
	} else {
		this.Data["agentVersion"] = ""
		this.Data["agentSpeed"] = 0
		this.Data["agentIsWaiting"] = false
	}
	this.Data["agent"] = agent

	// 分组
	groupNames := []string{}
	config := agents.SharedGroupConfig()
	for _, groupId := range agent.GroupIds {
		group := config.FindGroup(groupId)
		if group == nil {
			continue
		}
		groupNames = append(groupNames, group.Name)
	}
	this.Data["groupNames"] = groupNames

	this.Show()
}
