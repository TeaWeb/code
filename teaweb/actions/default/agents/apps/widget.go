package apps

import (
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type WidgetAction actions.Action

// Widget
func (this *WidgetAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	agentutils.InitAppData(this, params.AgentId, params.AppId, "widget")

	this.Show()
}
