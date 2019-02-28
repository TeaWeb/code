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
	one, err := query.Find()
	if err != nil {
		this.Fail(err.Error())
	}
	if one != nil {
		accessLog := one.(*tealogs.AccessLog)
		this.Data["headers"] = accessLog.SentHeader
		this.Data["body"] = string(accessLog.ResponseBodyData)
	} else {
		this.Data["headers"] = map[string][]string{}
		this.Data["body"] = ""
	}

	this.Success()
}
