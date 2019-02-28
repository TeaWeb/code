package index

import (
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	this.RedirectURL("/dashboard")
}
