package log

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"net/http"
	"time"
)

type GetAction actions.Action

// 获取日志
func (this *GetAction) Run(params struct {
	Server       string
	FromId       string
	Size         int64 `default:"10"`
	BodyFetching bool
}) {
	serverId := ""
	if len(params.Server) > 0 {
		server, err := teaconfigs.NewServerConfigFromFile(params.Server)
		if err != nil {
			this.Fail("发生错误：" + err.Error())
		}
		serverId = server.Id
	}

	requestBodyFetching = params.BodyFetching
	requestBodyTime = time.Now()

	logger := tealogs.SharedLogger()
	accessLogs := lists.NewList(logger.ReadNewLogs(serverId, params.FromId, params.Size))
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
		}
	})

	this.Data["logs"] = result.Slice

	this.Success()
}
