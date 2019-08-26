package log

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teadb"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type HistoryAction actions.Action

// 历史日志
func (this *HistoryAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	proxyutils.AddServerMenu(this)

	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err := teamongo.Test()
	mongoAvailable := true
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接MongoDB"
		mongoAvailable = false
	}

	this.Data["server"] = maps.Map{
		"id": params.ServerId,
	}

	// 列出最近30天的日志
	days := []maps.Map{}
	if mongoAvailable {
		for i := 0; i < 60; i++ {
			day := timeutil.Format("Ymd", time.Now().Add(time.Duration(-i*24)*time.Hour))

			b, err := teadb.AccessLogDAO().HasAccessLog(day, server.Id)
			if err != nil {
				logs.Error(err)
			}
			if b {
				days = append(days, maps.Map{
					"day": day,
					"has": true,
				})
			} else {
				days = append(days, maps.Map{
					"day": day,
					"has": false,
				})
			}
		}
	}

	this.Data["days"] = days
	this.Data["today"] = timeutil.Format("Ymd")

	this.Show()
}
