package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaplugins"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	this.Data["teaMenu"] = "dashboard"

	groups := teaplugins.DashboardGroups()
	for _, group := range groups {
		group.ForceReload()
	}

	this.Show()
}
