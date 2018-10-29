package login

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	username := this.Session().GetString("username")
	this.Data["username"] = username
	this.Data["passwordMask"] = ""

	config := configs.SharedAdminConfig()
	var found = false
	for _, user := range config.Users {
		if user.Username == username {
			this.Data["passwordMask"] = strings.Repeat("*", len(user.Password))
			found = true
		}
	}

	if !found {
		this.RedirectURL("/logout")
		return
	}

	this.Show()
}
