package agents

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAgentsAction actions.Action

// 删除一组Agents
func (this *DeleteAgentsAction) Run(params struct {
	AgentIds []string
}) {
	for _, agentId := range params.AgentIds {
		// 跳过本地主机
		if agentId == "local" {
			continue
		}
		agent := agents.NewAgentConfigFromId(agentId)
		if agent == nil {
			continue
		}

		agentList, err := agents.SharedAgentList()
		if err != nil {
			logs.Error(err)
			continue
		}
		agentList.RemoveAgent(agent.Filename())
		err = agentList.Save()
		if err != nil {
			logs.Error(err)
			continue
		}

		err = agent.Delete()
		if err != nil {
			logs.Error(err)
			continue
		}

		// 通知更新
		agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("REMOVE_AGENT", maps.Map{}))
	}

	this.Success()
}
