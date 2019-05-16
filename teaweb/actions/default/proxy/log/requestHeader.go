package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
)

type RequestHeaderAction actions.Action

// 请求Header
func (this *RequestHeaderAction) Run(params struct {
	LogId string
	Day   string
}) {
	query := teamongo.NewQuery("logs."+params.Day, new(tealogs.AccessLog))
	query.Id(params.LogId)
	query.Result("header", "requestData")
	one, err := query.Find()
	if err != nil {
		this.Fail(err.Error())
	}
	if one != nil {
		accessLog := one.(*tealogs.AccessLog)
		this.Data["headers"] = accessLog.Header
		this.Data["body"] = string(accessLog.RequestData)
	} else {
		this.Data["headers"] = map[string][]string{}
		this.Data["body"] = ""
	}

	this.Success()
}
