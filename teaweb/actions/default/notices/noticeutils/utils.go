package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// 删除Agent相关通知
func DeleteNoticesForAgent(agentId string) error {
	return NewNoticeQuery().
		Agent(&notices.AgentCond{
			AgentId: agentId,
		}).
		Action(NoticeQueryActionCount).
		Delete()
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
func CountReceivedNotices(receiverId string, cond map[string]interface{}, minutes int) int {
	if len(receiverId) == 0 {
		return 0
	}
	if minutes <= 0 {
		return 0
	}
	query := NewNoticeQuery().
		Attr("receivers", receiverId).
		Gte("timestamp", time.Now().Unix()-int64(minutes*60))
	if len(cond) > 0 {
		for k, v := range cond {
			query.Attr(k, v)
		}
	}
	c, err := query.
		Count()
	if err != nil {
		logs.Error(err)
	}
	return types.Int(c)
}

// 更改某个通知的接收人
func UpdateNoticeReceivers(id primitive.ObjectID, receiverIds []string) {
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

// 计算同样的消息数量
func ExistNoticesWithHash(hash string, cond map[string]interface{}, duration time.Duration) bool {
	query := NewNoticeQuery()
	query.Attr("messageHash", hash)
	for k, v := range cond {
		query.Attr(k, v)
	}
	query.Gt("timestamp", float64(time.Now().Unix())-duration.Seconds())
	query.Desc("_id")
	notice, err := query.Find()
	if err != nil {
		logs.Error(err)
		return false
	}
	if notice == nil {
		return false
	}

	// 中间是否有success级别的
	query2 := NewNoticeQuery()
	for k, v := range cond {
		query2.Attr(k, v)
	}
	if len(notice.Proxy.ServerId) > 0 {
		query2.Attr("proxy.level", notices.NoticeLevelSuccess)
		query2.Gt("_id", notice.Id)
	} else if len(notice.Agent.AgentId) > 0 {
		query2.Attr("agent.level", notices.NoticeLevelSuccess)
		query2.Gt("_id", notice.Id)
	}
	result, err := query2.Find()
	if err != nil {
		logs.Error(err)
	}
	return result == nil
}
