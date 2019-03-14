package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemValuesAction actions.Action

// 监控项数据展示
func (this *ItemValuesAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	Level   int
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")
	item := app.FindItem(params.ItemId)

	if item == nil {
		this.Fail("找不到要查看的Item")
	}

	this.Data["item"] = item
	this.Data["levels"] = notices.AllNoticeLevels()
	this.Data["selectedLevel"] = params.Level

	this.Show()
}

// 获取监控项数据
func (this *ItemValuesAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string
	LastId  string
	Level   notices.NoticeLevel
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	query := teamongo.NewAgentValueQuery()
	query.Agent(params.AgentId)
	query.App(params.AppId)
	query.Item(params.ItemId)
	query.Offset(0)
	query.Limit(100)
	query.Desc("_id")
	query.Action(teamongo.ValueQueryActionFindAll)

	if params.Level > 0 {
		if params.Level == notices.NoticeLevelInfo {
			query.Attr("noticeLevel", []interface{}{notices.NoticeLevelInfo, notices.NoticeLevelNone})
		} else {
			query.Attr("noticeLevel", params.Level)
		}
	}

	if len(params.LastId) > 0 {
		lastObjectId, err := primitive.ObjectIDFromHex(params.LastId)
		if err != nil {
			logs.Error(err)
		} else {
			query.Gt("_id", lastObjectId)
		}
	}

	ones, err := query.Execute()
	if err != nil {
		this.Fail("查询失败：" + err.Error())
	}

	this.Data["values"] = lists.Map(ones, func(k int, v interface{}) interface{} {
		value := v.(*agents.Value)
		return maps.Map{
			"id":          value.Id.Hex(),
			"timestamp":   value.Timestamp,
			"timeFormat":  value.TimeFormat,
			"value":       value.Value,
			"error":       value.Error,
			"noticeLevel": notices.FindNoticeLevel(value.NoticeLevel),
			"threshold":   value.Threshold,
		}
	})
	this.Success()
}
