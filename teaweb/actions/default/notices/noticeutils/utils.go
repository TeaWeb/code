package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"time"
)

// 获取某个Agent的未读通知数
func CountUnreadNoticesForAgent(agentId string) int {
	countNotices, err := NewNoticeQuery().
		Agent(&notices.AgentCond{
			AgentId: agentId,
		}).
		Attr("isRead", false).
		Action(NoticeQueryActionCount).
		Execute()
	if err != nil {
		return 0
	} else {
		return types.Int(countNotices)
	}
}

// 获取所有未读通知数
func CountUnreadNotices() int {
	countNotices, err := NewNoticeQuery().
		Attr("isRead", false).
		Action(NoticeQueryActionCount).
		Execute()
	if err != nil {
		return 0
	} else {
		return types.Int(countNotices)
	}
}

// 获取某个Agent已读通知数
func CountReadNoticesForAgent(agentId string) int {
	countNotices, err := NewNoticeQuery().
		Agent(&notices.AgentCond{
			AgentId: agentId,
		}).
		Attr("isRead", true).
		Action(NoticeQueryActionCount).
		Execute()
	if err != nil {
		return 0
	} else {
		return types.Int(countNotices)
	}
}

// 获取所有已读通知数
func CountReadNotices() int {
	countNotices, err := NewNoticeQuery().
		Attr("isRead", true).
		Action(NoticeQueryActionCount).
		Execute()
	if err != nil {
		return 0
	} else {
		return types.Int(countNotices)
	}
}

// 获取某个接收人在某个时间段内接收的通知数
func CountReceivedNotices(receiverId string, minutes int) int {
	if len(receiverId) == 0 {
		return 0
	}
	if minutes <= 0 {
		return 0
	}
	c, err := NewNoticeQuery().
		Attr("receivers", receiverId).
		Gte("timestamp", time.Now().Unix()-int64(minutes*60)).
		Count()
	if err != nil {
		logs.Error(err)
	}
	return types.Int(c)
}

// 更改某个通知的接收人
func UpdateNoticeReceivers(id objectid.ObjectID, receiverIds []string) {
	err := NewNoticeQuery().
		Attr("_id", id).
		Update(maps.Map{
			"$set": maps.Map{
				"isNotified": true,
				"receivers":  receiverIds,
			},
		})
	if err != nil {
		logs.Error(err)
	}
}
