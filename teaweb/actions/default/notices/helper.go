package notices

import (
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/TeaWeb/code/teaweb/utils"
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

	// 操作按钮
	menuGroup := utils.NewMenuGroup()
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000
		menu.Add("通知", "", "/notices", true)
	}

	menuGroup.Sort()
	utils.SetSubMenu(action, menuGroup)
}
