package agents

import (
	"context"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"time"
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

		// 删除通知
		err := noticeutils.DeleteNoticesForAgent(agent.Id)
		if err != nil {
			this.Fail("通知删除失败：" + err.Error())
		}

		// 删除数值记录
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err = teamongo.FindCollection("values.agent." + agent.Id).Drop(ctx)
		if err != nil {
			this.Fail("数值记录删除失败：" + err.Error())
		}

		// 从列表删除
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

		// 删除通知
		err = noticeutils.DeleteNoticesForAgent(agent.Id)
		if err != nil {
			logs.Error(err)
		}

		// 通知更新
		agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("REMOVE_AGENT", maps.Map{}))
	}

	this.Success()
}
