package log

import (
	"fmt"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/time"
	"net/http"
	"time"
)

type ListAction actions.Action

// 获取日志
func (this *ListAction) Run(params struct {
	ServerId     string
	FromId       string
	Size         int64 `default:"10"`
	BodyFetching bool
	LogType      string
}) {

	if params.Size < 1 {
		params.Size = 20
	}

	serverId := params.ServerId

	requestBodyFetching = params.BodyFetching
	requestBodyTime = time.Now()

	shouldReverse := true
	query := teamongo.NewQuery("logs."+timeutil.Format("Ymd"), new(tealogs.AccessLog))
	query.Attr("serverId", serverId)
	if len(params.FromId) > 0 {
		query.Gt("_id", params.FromId)
		query.AscPk()
	} else {
		query.DescPk()
		shouldReverse = false
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
	query.Limit(params.Size)
	ones, err := query.FindAll()

	if err != nil {
		logs.Error(err)
		this.Data["logs"] = []interface{}{}
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
				"backendId":      accessLog.BackendId,
				"locationId":     accessLog.LocationId,
				"rewriteId":      accessLog.RewriteId,
				"fastcgiId":      accessLog.FastcgiId,
				"attrs":          accessLog.Attrs,
			}
		})

		if shouldReverse {
			lists.Reverse(result)
		}
		this.Data["logs"] = result
	}

	this.Success()
}
