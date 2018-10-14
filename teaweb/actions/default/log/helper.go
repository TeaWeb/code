package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "log"
	action.Data["teaTabbar"] = []maps.Map{
		{
			"name":    "访问日志",
			"subName": "",
			"url":     "/log",
			"active":  action.Spec.HasClassPrefix("log."),
		},
	}
}
