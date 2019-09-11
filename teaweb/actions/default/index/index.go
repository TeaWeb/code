package index

import (
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	// 是否已初始化
	config := db.SharedDBConfig()
	if !config.IsInitialized {
		this.RedirectURL("/install")
		return
	}

	this.RedirectURL("/dashboard")
}
