package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
)

type ResponseHeaderAction actions.Action

// 响应Header
func (this *ResponseHeaderAction) Run(params struct {
	LogId string
}) {
	accessLog := tealogs.SharedLogger().FindLog(params.LogId)
	if accessLog != nil {
		this.Data["headers"] = accessLog.SentHeader
	} else {
		this.Data["headers"] = map[string][]string{}
	}

	this.Data["body"] = string(accessLog.ResponseBodyData)

	this.Success()
}
