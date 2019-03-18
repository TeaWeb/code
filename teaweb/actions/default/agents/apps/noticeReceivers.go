package apps

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type NoticeReceiversAction actions.Action

// 通知接收人设置
func (this *NoticeReceiversAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "noticeSetting")
	if app == nil {
		this.Fail("找不到要操作的App")
	}

	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		level := v.(maps.Map)
		code := level["code"].(notices.NoticeLevel)
		receivers, found := app.NoticeSetting[code]
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
