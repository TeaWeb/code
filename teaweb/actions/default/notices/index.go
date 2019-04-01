package notices

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"math"
	"time"
)

type IndexAction actions.Action

// 通知
func (this *IndexAction) Run(params struct {
	Read int
	Page int
}) {
	this.Data["isRead"] = params.Read > 0

	count := 0
	countUnread := noticeutils.CountUnreadNotices()
	if params.Read == 0 {
		count = countUnread
	} else {
		count = noticeutils.CountReadNotices()
	}

	this.Data["countUnread"] = countUnread
	this.Data["count"] = count
	this.Data["soundOn"] = notices.SharedNoticeSetting().SoundOn

	// 分页
	if params.Page < 1 {
		params.Page = 1
	}
	pageSize := 10
	this.Data["page"] = params.Page
	if count > 0 {
		this.Data["countPages"] = int(math.Ceil(float64(count) / float64(pageSize)))
	} else {
		this.Data["countPages"] = 0
	}

	// 读取数据
	ones, err := noticeutils.NewNoticeQuery().
		Attr("isRead", params.Read == 1).
		Offset(int64((params.Page - 1) * pageSize)).
		Limit(int64(pageSize)).
		Desc("_id").
		FindAll()
	if err != nil {
		logs.Error(err)
		this.Data["notices"] = []maps.Map{}
	} else {
		this.Data["notices"] = lists.Map(ones, func(k int, v interface{}) interface{} {
			notice := v.(*notices.Notice)
			isAgent := len(notice.Agent.AgentId) > 0
			m := maps.Map{
				"id":       notice.Id,
				"isAgent":  isAgent,
				"isRead":   notice.IsRead,
				"message":  notice.Message,
				"datetime": timeutil.Format("Y-m-d H:i:s", time.Unix(notice.Timestamp, 0)),
			}

			// Agent
			if isAgent {
				m["level"] = notices.FindNoticeLevel(notice.Agent.Level)

				links := []maps.Map{}
				agent := agents.NewAgentConfigFromId(notice.Agent.AgentId)
				if agent != nil {
					links = append(links, maps.Map{
						"name": agent.Name,
						"url":  "/agents/board?agentId=" + agent.Id,
					})

					app := agent.FindApp(notice.Agent.AppId)
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

				m["links"] = links
			}

			return m
		})
	}

	this.Show()
}
