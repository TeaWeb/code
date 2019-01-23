package apps

import (
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type RunTaskAction actions.Action

// 运行一次任务
func (this *RunTaskAction) Run(params struct {
	AgentId string
	TaskId  string
}) {
	agentutils.PostAgentEvent(params.AgentId, &agentutils.Event{
		Name: "RUN_TASK",
		Data: maps.Map{
			"taskId": params.TaskId,
		},
	})

	this.Success()
}
