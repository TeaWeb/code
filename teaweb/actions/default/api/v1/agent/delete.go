package agent

import (
	"github.com/TeaWeb/code/teaweb/actions/default/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除Agent
func (this *DeleteAction) RunGet(params struct {
	AgentId string
}) {
	defer apiutils.Recover(this, false)

	act := new(agents.DeleteAction)
	act.Request = this.Request
	act.ResponseWriter = new(actions.TestingResponseWriter)
	act.RunPost(struct {
		AgentId string
	}{
		AgentId: params.AgentId,
	})
}
