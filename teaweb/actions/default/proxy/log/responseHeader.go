package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
)

type ResponseHeaderAction actions.Action

// 响应Header
func (this *ResponseHeaderAction) Run(params struct {
	LogId string
	Day   string
}) {
	query := teamongo.NewQuery("logs."+params.Day, new(tealogs.AccessLog))
	query.Id(params.LogId)
	query.Result("sentHeader", "responseBodyData")
	one, err := query.Find()
	if err != nil {
		this.Fail(err.Error())
	}
	if one != nil {
		accessLog := one.(*tealogs.AccessLog)
		this.Data["headers"] = accessLog.SentHeader
		this.Data["hasBody"] = len(accessLog.ResponseBodyData) > 0
	} else {
		this.Data["headers"] = map[string][]string{}
		this.Data["hasBody"] = false
	}

	this.Success()
}
