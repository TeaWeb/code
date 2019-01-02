package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
	"time"
)

type CookiesAction actions.Action

func (this *CookiesAction) Run(params struct {
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
		this.Data["cookies"] = accessLog.Cookie
		this.Data["count"] = len(accessLog.Cookie)
	} else {
		this.Data["cookies"] = map[string]string{}
		this.Data["count"] = 0
	}

	this.Success()
}
