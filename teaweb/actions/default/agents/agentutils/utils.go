package agentutils

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/maps"
)

// 查找通知相关的Link
func FindNoticeLinks(notice *notices.Notice) (links []maps.Map) {
	if len(notice.Agent.AgentId) > 0 {
		agent := agents.NewAgentConfigFromId(notice.Agent.AgentId)
		if agent != nil {
			links = append(links, maps.Map{
				"name": agent.Name,
				"url":  "/agents/board?agentId=" + agent.Id,
			})

			app := FindAgentApp(agent, notice.Agent.AppId)
			if app != nil {
				links = append(links, maps.Map{
					"name": app.Name,
					"url":  "/agents/apps/detail?agentId=" + agent.Id + "&appId=" + app.Id,
				})

				// item
				if len(notice.Agent.ItemId) > 0 {
					item := app.FindItem(notice.Agent.ItemId)
					if item != nil {
						links = append(links, maps.Map{
							"name": item.Name,
							"url":  "/agents/apps/itemDetail?agentId=" + agent.Id + "&appId=" + app.Id + "&itemId=" + item.Id,
						})
					}
				}

				// task
				if len(notice.Agent.TaskId) > 0 {
					task := app.FindTask(notice.Agent.TaskId)
					if task != nil {
						links = append(links, maps.Map{
							"name": task.Name,
							"url":  "/agents/apps/itemDetail?agentId=" + agent.Id + "&appId=" + app.Id + "&taskId=" + task.Id,
						})
					}
				}
			}
		}
	}
	return
}
