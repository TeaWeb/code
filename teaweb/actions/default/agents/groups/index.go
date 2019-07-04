package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
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

	countDefaultReceivers := 0
	groups := []maps.Map{}
	for _, group := range agents.SharedGroupConfig().Groups {
		countAgents := 0
		for _, agent := range allAgents {
			if lists.ContainsString(agent.GroupIds, group.Id) {
				countAgents++
			}
		}

		countReceivers := 0
		for _, receivers := range group.NoticeSetting {
			countReceivers += len(receivers)
		}

		if len(group.Id) == 0 {
			countDefaultReceivers = countReceivers
			continue
		} else {
			groups = append(groups, maps.Map{
				"id":             group.Id,
				"name":           group.Name,
				"on":             group.On,
				"countAgents":    countAgents,
				"countReceivers": countReceivers,
				"canDelete":      true,
			})
		}
	}

	// 默认分组
	countDefaultAgents := 0
	for _, agent := range allAgents {
		if len(agent.GroupIds) == 0 && len(agent.Id) > 0 {
			countDefaultAgents++
		}
	}

	defaultGroupObject := agents.LoadDefaultGroup()
	defaultGroup := maps.Map{
		"id":             "",
		"name":           "[" + defaultGroupObject.Name + "]",
		"on":             true,
		"countAgents":    countDefaultAgents,
		"countReceivers": countDefaultReceivers,
		"canDelete":      false,
	}

	this.Data["groups"] = append([]maps.Map{defaultGroup}, groups ...)
	this.Data["noticeLevels"] = notices.AllNoticeLevels()

	this.Show()
}
