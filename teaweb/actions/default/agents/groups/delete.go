package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除分组
func (this *DeleteAction) Run(params struct {
	GroupId string
}) {
	// 删除agent中的groupId
	agentList, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	for _, agent := range agentList.FindAllAgents() {
		agent.RemoveGroup(params.GroupId)
		err = agent.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	config := agents.SharedGroupConfig()
	config.RemoveGroup(params.GroupId)
	err = config.Save()
	if err != nil {
		this.Fail("保存失败： " + err.Error())
	}

	this.Success()
}
