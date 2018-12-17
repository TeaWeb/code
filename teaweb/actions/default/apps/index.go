package apps

import (
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 所有Apps
func (this *IndexAction) Run(params struct{}) {
	this.RedirectURL("/apps/all")
	this.Show()
}
