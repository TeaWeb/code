package notices

import (
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	action.Data["teaMenu"] = "notices"

	if action.Request.Method == http.MethodGet {
		if !action.HasPrefix("/notices/badge") {
			action.Data["countUnread"] = noticeutils.CountUnreadNotices()
		}
	}
}
