package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaplugins"
	"time"
)

type WidgetsAction actions.Action

func (this *WidgetsAction) Run() {
	groups := teaplugins.DashboardGroups()
	for _, group := range groups {
		group.Reload()
	}

	// 暂停，等待reload()任务完成
	time.Sleep(100 * time.Millisecond)

	this.Data["widgetGroups"] = groups

	this.Success()
}
