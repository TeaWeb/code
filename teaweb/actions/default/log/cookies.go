package log

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
)

type CookiesAction actions.Action

func (this *CookiesAction) Run(params struct {
	LogId string
}) {
	accessLog := tealogs.SharedLogger().FindLog(params.LogId)
	if accessLog != nil {
		this.Data["cookies"] = accessLog.Cookie
		this.Data["count"] = len(accessLog.Cookie)
	} else {
		this.Data["cookies"] = map[string]string{}
		this.Data["count"] = 0
	}

	this.Success()
}
