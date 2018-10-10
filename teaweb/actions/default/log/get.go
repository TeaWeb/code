package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/tealogs"
	"time"
	"github.com/iwind/TeaGo/lists"
)

type GetAction actions.Action

func (this *GetAction) Run(params struct {
	FromId string
	Size   int64 `default:"10"`
}) {
	logger := tealogs.SharedLogger()
	accessLogs := lists.NewList(logger.ReadNewLogs(params.FromId, params.Size))
	result := accessLogs.Map(func(k int, v interface{}) interface{} {
		accessLog := v.(tealogs.AccessLog)
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
			"statusMessage":  accessLog.StatusMessage,
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
		}
	})

	this.Data["logs"] = result.Slice

	fromTime := time.Now().Add(-24 * time.Hour)
	toTime := time.Now()

	countSuccess := logger.CountSuccessLogs(fromTime.Unix(), toTime.Unix())
	countFail := logger.CountFailLogs(fromTime.Unix(), toTime.Unix())
	total := countSuccess + countFail
	this.Data["countSuccess"] = countSuccess
	this.Data["countFail"] = countFail
	this.Data["total"] = total
	this.Data["qps"] = logger.QPS()

	this.Success()
}
