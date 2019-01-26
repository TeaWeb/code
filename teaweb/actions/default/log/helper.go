package log

import (
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	tabbar := utils.NewTabbar()
	tabbar.Add("系统日志", "", "/log/runtime", "", action.HasPrefix("/log/runtime"))
	tabbar.Add("操作日志", "", "/log/audit", "", action.HasPrefix("/log/audit"))

	action.Data["teaTabbar"] = tabbar.Items()
	action.Data["teaMenu"] = "log.runtime"
}
