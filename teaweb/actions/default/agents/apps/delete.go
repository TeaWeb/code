package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAction actions.Action

// 删除App
func (this *DeleteAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要操作的Agent")
	}

	// 删除图表
	board := agents.NewAgentBoard(agent.Id)
	if board != nil {
		board.RemoveApp(params.AppId)
		err := board.Save()
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}

	// 删除App
	agent.RemoveApp(params.AppId)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("REMOVE_APP", maps.Map{
		"appId": params.AppId,
	}))

	this.Success()
}
