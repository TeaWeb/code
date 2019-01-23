package apps

import (
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type AddWidgetAction actions.Action

// 添加Widget
func (this *AddWidgetAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	agentutils.InitAppData(this, params.AgentId, params.AppId, "widget")

	this.Show()
}
