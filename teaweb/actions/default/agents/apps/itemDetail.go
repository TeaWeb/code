package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ItemDetailAction actions.Action

// 监控项详情
func (this *ItemDetailAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	this.Data["agent"] = maps.Map{
		"id": params.AgentId,
	}

	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")
	item := app.FindItem(params.ItemId)

	if item == nil {
		this.Fail("找不到要查看的Item")
	}

	this.Data["item"] = item

	this.Data["source"] = nil
	source := item.Source()
	if source != nil {
		this.Data["source"] = maps.Map{
			"summary":    agents.FindDataSource(source.Code()),
			"options":    source,
			"dataFormat": agents.FindSourceDataFormat(source.DataFormatCode()),
		}
	}

	this.Data["noticeLevels"] = notices.AllNoticeLevels()

	this.Show()
}
