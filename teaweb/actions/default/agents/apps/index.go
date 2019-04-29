package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板首页
func (this *IndexAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	// 用户自定义App
	this.Data["apps"] = lists.Map(agent.Apps, func(k int, v interface{}) interface{} {
		app := v.(*agents.AppConfig)

		// 最新一条数据
		level := notices.NoticeLevelNone
		for _, item := range app.Items {
			if !item.On {
				continue
			}
			value, err := teamongo.NewAgentValueQuery().
				Agent(agent.Id).
				App(app.Id).
				Item(item.Id).
				Desc("_id").
				Find()
			if err == nil && value != nil {
				if value.NoticeLevel == notices.NoticeLevelWarning || value.NoticeLevel == notices.NoticeLevelError && value.NoticeLevel > level {
					level = value.NoticeLevel
				}
			}
		}

		return maps.Map{
			"on":                app.On,
			"id":                app.Id,
			"name":              app.Name,
			"items":             app.Items,
			"bootingTasks":      app.FindBootingTasks(),
			"manualTasks":       app.FindManualTasks(),
			"schedulingTasks":   app.FindSchedulingTasks(),
			"isSharedWithGroup": app.IsSharedWithGroup,
			"isWarning":         level == notices.NoticeLevelWarning,
			"isError":           level == notices.NoticeLevelError,
		}
	})

	this.Show()
}
