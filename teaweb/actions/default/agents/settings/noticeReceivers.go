package settings

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type NoticeReceiversAction actions.Action

// 通知接收人设置
func (this *NoticeReceiversAction) Run(params struct {
	AgentId string
}) {
	this.Data["selectedTab"] = "noticeSetting"

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}
	this.Data["agent"] = agent

	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		level := v.(maps.Map)
		code := level["code"].(notices.NoticeLevel)
		receivers, found := agent.NoticeSetting[code]
		if found && len(receivers) > 0 {
			level["receivers"] = lists.Map(receivers, func(k int, v interface{}) interface{} {
				receiver := v.(*notices.NoticeReceiver)

				m := maps.Map{
					"name":      receiver.Name,
					"id":        receiver.Id,
					"user":      receiver.User,
					"mediaType": "",
				}

				// 媒介
				media := notices.SharedNoticeSetting().FindMedia(receiver.MediaId)
				if media != nil {
					m["mediaType"] = media.Name
				}

				return m
			})
		} else {
			level["receivers"] = []interface{}{}
		}
		return level
	})

	this.Show()
}
