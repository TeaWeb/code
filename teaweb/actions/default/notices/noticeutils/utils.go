package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/types"
)

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
