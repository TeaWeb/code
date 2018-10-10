package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaplugins"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	// 检查mongodb
	err := teamongo.Test()
	if err != nil {
		this.RedirectURL("/install/mongo")
		return
	}

	this.Data["teaMenu"] = "dashboard"

	groups := teaplugins.DashboardGroups()
	for _, group := range groups {
		group.ForceReload()
	}

	this.Show()
}
