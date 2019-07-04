package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type NoticeReceiversAction actions.Action

// 通知接收人
func (this *NoticeReceiversAction) Run(params struct {
	GroupId string
}) {
	group := agents.SharedGroupConfig().FindGroup(params.GroupId)
	if group == nil {
		if len(params.GroupId) == 0 {
			group = &agents.Group{
				Id:   "",
				Name: agents.LoadDefaultGroup().Name,
				On:   true,
			}
		} else {
			this.Fail("Group不存在")
		}
	}

	this.Data["group"] = group
	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		level := v.(maps.Map)
		code := level["code"].(notices.NoticeLevel)
		receivers, found := group.NoticeSetting[code]
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
