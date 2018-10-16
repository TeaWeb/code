package mongo

import "github.com/iwind/TeaGo/actions"

type InstallStatusAction actions.Action

func (this *InstallStatusAction) Run(params struct{}) {
	this.Data["status"] = installStatus
	this.Data["percent"] = installPercent

	this.Success()
}
