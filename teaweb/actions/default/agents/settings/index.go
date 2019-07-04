package settings

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
	this.Data["defaultGroupName"] = agents.LoadDefaultGroup().Name
	this.Data["selectedTab"] = "detail"

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	state, isWaiting := agentutils.CheckAgentIsWaiting(agent.Id)
	if isWaiting {
		this.Data["agentVersion"] = state.Version
		this.Data["agentSpeed"] = state.Speed
		this.Data["agentIP"] = state.IP
		this.Data["agentIsWaiting"] = true
	} else {
		this.Data["agentVersion"] = ""
		this.Data["agentSpeed"] = 0
		this.Data["agentIP"] = ""
		this.Data["agentIsWaiting"] = false
	}
	this.Data["agent"] = agent
	this.Data["isLocal"] = agent.IsLocal()

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
