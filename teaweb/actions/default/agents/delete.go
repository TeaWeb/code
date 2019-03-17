package agents

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId

	if params.AgentId == "local" {
		this.RedirectURL("/agents/board")
		return
	}

	this.Show()
}

// 提交
func (this *DeleteAction) RunPost(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("要删除的主机不存在")
	}

	agentList, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}
	agentList.RemoveAgent(agent.Filename())
	err = agentList.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = agent.Delete()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	// 删除通知
	err = noticeutils.DeleteNoticesForAgent(agent.Id)
	if err != nil {
		logs.Error(err)
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("REMOVE_AGENT", maps.Map{}))

	this.Success()
}
