package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 详情
func (this *DetailAction) Run(params struct {
	GroupId string
}) {
	if len(params.GroupId) == 0 {
		this.Data["group"] = agents.LoadDefaultGroup()
	} else {
		group := agents.SharedGroupConfig().FindGroup(params.GroupId)
		if group == nil {
			this.Fail("找不到Group")
		}
		this.Data["group"] = group
	}

	// Agents列表
	groupAgents := []maps.Map{}
	agentList, err := agents.SharedAgentList()
	if err != nil {
		logs.Error(err)
	} else {
		for _, a := range agentList.FindAllAgents() {
			if (len(a.GroupIds) == 0 && len(params.GroupId) == 0) || lists.ContainsString(a.GroupIds, params.GroupId) {
				_, isWaiting := agentutils.CheckAgentIsWaiting(a.Id)
				groupAgents = append(groupAgents, maps.Map{
					"on":        a.On,
					"id":        a.Id,
					"name":      a.Name,
					"host":      a.Host,
					"isWaiting": isWaiting,
				})
			}
		}
	}

	this.Data["agents"] = groupAgents

	this.Show()
}
