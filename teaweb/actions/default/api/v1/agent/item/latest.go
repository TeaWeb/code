package item

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type LatestAction actions.Action

// 获取最后一次数据
func (this *LatestAction) RunGet(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "agent not found")
		return
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		apiutils.Fail(this, "app not found")
		return
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		apiutils.Fail(this, "item not found")
		return
	}

	value, err := teamongo.NewAgentValueQuery().
		Agent(params.AgentId).
		App(params.AppId).
		Item(item.Id).
		Desc("_id").
		Find()
	if err != nil {
		apiutils.Fail(this, "no value yet")
		return
	}

	apiutils.Success(this, value.Value)
}
