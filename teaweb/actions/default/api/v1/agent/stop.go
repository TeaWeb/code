package agent

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type StopAction actions.Action

// 关闭Agent
func (this *StopAction) RunGet(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "not found")
		return
	}

	agent.On = false
	err := agent.Save()
	if err != nil {
		apiutils.Fail(this, "保存失败:"+err.Error())
		return
	}

	apiutils.SuccessOK(this)
}
