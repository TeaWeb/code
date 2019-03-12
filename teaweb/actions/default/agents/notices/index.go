package notices

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
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

// 通知首页
func (this *IndexAction) Run(params struct {
	AgentId string
	Read    int
	Page    int
}) {
	this.Data["agentId"] = params.AgentId
	this.Data["isRead"] = params.Read > 0

	count := 0
	countUnread := noticeutils.CountUnreadNoticesForAgent(params.AgentId)
	if params.Read == 0 {
		count = countUnread
	} else {
		count = noticeutils.CountReadNoticesForAgent(params.AgentId)
	}

	this.Data["countUnread"] = countUnread
	this.Data["count"] = count

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
		Agent(&notices.AgentCond{
			AgentId: params.AgentId,
		}).
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
			if len(notice.Agent.Threshold) > 0 {
				notice.Message += " [触发阈值：" + notice.Agent.Threshold + "]"
			}
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
				m["links"] = agentutils.FindNoticeLinks(notice)
			}

			return m
		})
	}

	this.Show()
}
