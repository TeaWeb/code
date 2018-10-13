package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
)

type RequestHeaderAction actions.Action

func (this *RequestHeaderAction) Run(params struct {
	LogId string
}) {
	accessLog := tealogs.SharedLogger().FindLog(params.LogId)
	if accessLog != nil {
		this.Data["headers"] = accessLog.Header
	} else {
		this.Data["headers"] = map[string][]string{}
	}

	this.Success()
}
