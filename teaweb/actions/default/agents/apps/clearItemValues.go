package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
)

type ClearItemValuesAction actions.Action

// 清除数值记录
func (this *ClearItemValuesAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	query := teamongo.NewAgentValueQuery()
	query.Agent(agent.Id)
	query.Attr("appId", app.Id)
	query.Attr("itemId", item.Id)
	err := query.Delete()
	if err != nil {
		this.Fail("清除失败：" + err.Error())
	}

	this.Success()
}
