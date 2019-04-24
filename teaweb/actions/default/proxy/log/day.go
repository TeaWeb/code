package log

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"net/http"
	"regexp"
	"time"
)

type DayAction actions.Action

// 某天的日志
func (this *DayAction) Run(params struct {
	ServerId string
	Day      string
	LogType  string
	FromId   string
	Page     int
	Size     int64
	SearchIP string
}) {
	serverId := params.ServerId
	server := teaconfigs.NewServerConfigFromId(serverId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 {
		params.Size = 20
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}
	this.Data["searchIP"] = params.SearchIP

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

	this.Data["day"] = params.Day
	this.Data["isHistory"] = regexp.MustCompile("^\\d+$").MatchString(params.Day)
	this.Data["logType"] = params.LogType
	this.Data["logs"] = []interface{}{}
	this.Data["fromId"] = params.FromId
	this.Data["hasNext"] = false
	this.Data["page"] = params.Page

	// 日志列表
	if mongoAvailable {
		realDay := ""
		if regexp.MustCompile("^\\d+$").MatchString(params.Day) {
			realDay = params.Day
		} else if params.Day == "today" {
			realDay = timeutil.Format("Ymd")
		} else if params.Day == "yesterday" {
			realDay = timeutil.Format("Ymd", time.Now().Add(-24*time.Hour))
		} else {
			realDay = timeutil.Format("Ymd")
		}

		query := teamongo.NewQuery("logs."+realDay, new(tealogs.AccessLog))
		query.Attr("serverId", serverId)
		if len(params.FromId) > 0 {
			query.Lte("_id", params.FromId)
		}
		if params.LogType == "errorLog" {
			query.Or([]map[string]interface{}{
				{
					"hasErrors": true,
				},
				{
					"status": map[string]interface{}{
						"$gte": 400,
					},
				},
			} ...)
		}
		if len(params.SearchIP) > 0 {
			query.Attr("remoteAddr", params.SearchIP)
		}
		query.Offset(params.Size * int64(params.Page-1))
		query.Limit(params.Size)
		query.DescPk()
		ones, err := query.FindAll()
		if err != nil {
			this.Data["mongoError"] = "MongoDB查询错误：" + err.Error()
		} else {
			result := lists.Map(ones, func(k int, v interface{}) interface{} {
				accessLog := v.(*tealogs.AccessLog)
				return map[string]interface{}{
					"id":             accessLog.Id.Hex(),
					"requestTime":    accessLog.RequestTime,
					"request":        accessLog.Request,
					"requestURI":     accessLog.RequestURI,
					"requestMethod":  accessLog.RequestMethod,
					"remoteAddr":     accessLog.RemoteAddr,
					"remotePort":     accessLog.RemotePort,
					"userAgent":      accessLog.UserAgent,
					"host":           accessLog.Host,
					"status":         accessLog.Status,
					"statusMessage":  fmt.Sprintf("%d", accessLog.Status) + " " + http.StatusText(accessLog.Status),
					"timeISO8601":    accessLog.TimeISO8601,
					"timeLocal":      accessLog.TimeLocal,
					"requestScheme":  accessLog.Scheme,
					"proto":          accessLog.Proto,
					"contentType":    accessLog.SentContentType(),
					"bytesSent":      accessLog.BytesSent,
					"backendAddress": accessLog.BackendAddress,
					"fastcgiAddress": accessLog.FastcgiAddress,
					"extend":         accessLog.Extend,
					"referer":        accessLog.Referer,
					"upgrade":        accessLog.GetHeader("Upgrade"),
					"day":            timeutil.Format("Ymd", accessLog.Time()),
					"errors":         accessLog.Errors,
					"attrs":          accessLog.Attrs,
				}
			})

			this.Data["logs"] = result

			if len(result) > 0 {
				if len(params.FromId) == 0 {
					fromId := ones[0].(*tealogs.AccessLog).Id.Hex()
					this.Data["fromId"] = fromId
				}

				{
					nextId := ones[len(ones)-1].(*tealogs.AccessLog).Id.Hex()

					query := teamongo.NewQuery("logs."+realDay, new(tealogs.AccessLog))
					query.Attr("serverId", serverId)
					query.Lt("_id", nextId)
					if params.LogType == "errorLog" {
						query.Or([]map[string]interface{}{
							{
								"hasErrors": true,
							},
							{
								"status": map[string]interface{}{
									"$gte": 400,
								},
							},
						} ...)
					}
					if len(params.SearchIP) > 0 {
						query.Attr("remoteAddr", params.SearchIP)
					}
					v, err := query.Find()
					if err != nil {
						logs.Error(err)
					} else if v != nil {
						this.Data["hasNext"] = true
					}
				}
			}
		}
	}

	this.Show()
}
