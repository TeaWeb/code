package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
	"time"
)

type RequestHeaderAction actions.Action

// 请求Header
func (this *RequestHeaderAction) Run(params struct {
	LogId string
}) {
	query := tealogs.NewQuery()
	query.From(time.Now())
	query.Id(params.LogId)
	accessLog, err := query.Find()
	if err != nil {
		this.Fail(err.Error())
	}

	if accessLog != nil {
		this.Data["headers"] = accessLog.Header
	} else {
		this.Data["headers"] = map[string][]string{}
	}

	this.Data["body"] = string(accessLog.RequestData)

	this.Success()
}
