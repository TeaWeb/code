package agents

import (
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
}

// 筛选Action调用之前
func (this *Helper) BeforeAction(action actions.ActionWrapper) {
	agentutils.AddTabbar(action)
}
