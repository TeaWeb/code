package app

import (
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/iwind/TeaGo/actions"
)

type FavorAction actions.Action

func (this *FavorAction) Run(params struct {
	AppId string
}) {
	apputils.FavorApp(params.AppId)

	this.Success()
}
