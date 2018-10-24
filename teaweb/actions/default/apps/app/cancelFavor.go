package app

import (
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/iwind/TeaGo/actions"
)

type CancelFavorAction actions.Action

func (this *CancelFavorAction) Run(params struct {
	AppId string
}) {
	apputils.CancelFavorApp(params.AppId)

	this.Success()
}
