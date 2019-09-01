package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 分组管理
func (this *IndexAction) Run(params struct{}) {
	agentsList, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("ERROR:" + err.Error())
	}
	allAgents := agentsList.FindAllAgents()

	groups := []maps.Map{}
	for _, group := range agents.SharedGroupList().Groups {
		countAgents := 0
		for _, agent := range allAgents {
			if agent.BelongsToGroup(group.Id) {
				countAgents++
			}
		}

		countReceivers := 0
		for _, receivers := range group.NoticeSetting {
			countReceivers += len(receivers)
		}

		groups = append(groups, maps.Map{
			"id":             group.Id,
			"name":           group.Name,
			"on":             group.On,
			"countAgents":    countAgents,
			"countReceivers": countReceivers,
			"canDelete":      !group.IsDefault,
		})
	}

	this.Data["groups"] = groups
	this.Data["noticeLevels"] = notices.AllNoticeLevels()

	this.Show()
}
