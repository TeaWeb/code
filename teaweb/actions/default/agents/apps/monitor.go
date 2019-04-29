package apps

import (
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type MonitorAction actions.Action

// 监控
func (this *MonitorAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")
	err := app.Validate()
	if err != nil {
		logs.Error(err)
	}

	m := this.Data["app"].(maps.Map)
	m["items"] = app.Items

	this.Data["noticeLevels"] = notices.AllNoticeLevels()

	this.Show()
}

// 监控数据
func (this *MonitorAction) RunPost(params struct {
	AgentId string
	AppId   string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	this.Data["items"] = lists.Map(app.Items, func(k int, v interface{}) interface{} {
		item := v.(*agents.Item)

		latestValue := ""
		latestTime := ""
		latestLevel := notices.NoticeLevelNone

		value, err := teamongo.NewAgentValueQuery().
			Agent(params.AgentId).
			App(params.AppId).
			Item(item.Id).
			Desc("_id").
			Find()
		if err != nil {
			logs.Error(err)
		} else if value != nil {
			data, err := json.MarshalIndent(value.Value, "", "  ")
			if err != nil {
				logs.Error(err)
			} else {
				latestValue = string(data)
				latestTime = timeutil.Format("Y-m-d H:i:s", time.Unix(value.Timestamp, 0))
				latestLevel = value.NoticeLevel
			}
		}

		item.Validate()

		result := maps.Map{
			"id":          item.Id,
			"name":        item.Name,
			"on":          item.On,
			"interval":    item.Interval,
			"thresholds":  item.Thresholds,
			"latestValue": latestValue,
			"latestTime":  latestTime,
			"isWarning":   latestLevel == notices.NoticeLevelWarning,
			"isError":     latestLevel == notices.NoticeLevelError,
		}

		source := item.Source()
		if source != nil {
			result["sourceName"] = source.Name()
		} else {
			result["sourceName"] = "Agent"
		}

		return result
	})

	this.Success()
}
