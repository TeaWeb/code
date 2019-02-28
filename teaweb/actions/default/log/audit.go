package log

import (
	"github.com/TeaWeb/code/teaconfigs/audits"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"math"
	"time"
)

type AuditAction actions.Action

// 审计日志
func (this *AuditAction) Run(params struct {
	Read int
	Page int
}) {
	// 分页
	if params.Page < 1 {
		params.Page = 1
	}
	pageSize := 10
	this.Data["page"] = params.Page

	count, err := teamongo.NewAuditsQuery().Count()
	if err != nil {
		logs.Error(err)
	}

	if count > 0 {
		this.Data["countPages"] = int(math.Ceil(float64(count) / float64(pageSize)))
	} else {
		this.Data["countPages"] = 0
	}

	// 读取列表数据
	ones, err := teamongo.NewAuditsQuery().
		DescPk().
		Offset(int64(pageSize * (params.Page - 1))).
		Limit(int64(pageSize)).
		FindAll()
	if err != nil {
		this.Data["logs"] = []interface{}{}
	} else {
		this.Data["logs"] = lists.Map(ones, func(k int, v interface{}) interface{} {
			log := v.(*audits.Log)
			return maps.Map{
				"username":    log.Username,
				"action":      log.Action,
				"actionName":  log.ActionName(),
				"description": log.Description,
				"datetime":    timeutil.Format("Y-m-d H:i:s", time.Unix(log.Timestamp, 0)),
				"options":     log.Options,
			}
		})
	}

	this.Show()
}
