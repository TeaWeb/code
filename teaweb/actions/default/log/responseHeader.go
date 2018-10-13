package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
)

type ResponseHeaderAction actions.Action

func (this *ResponseHeaderAction) Run(params struct {
	LogId string
}) {
	accessLog := tealogs.SharedLogger().FindLog(params.LogId)
	if accessLog != nil {
		this.Data["headers"] = accessLog.SentHeader
	} else {
		this.Data["headers"] = map[string][]string{}
	}

	this.Success()
}
