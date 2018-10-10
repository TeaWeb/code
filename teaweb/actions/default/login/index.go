package login

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
)

type IndexAction actions.Action

var countLoginTries = 0

func (this *IndexAction) RunGet() {
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Username string
	Password string
	Must     *actions.Must
	Auth     *helpers.UserShouldAuth
}) {
	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	if countLoginTries >= 3 {
		this.Fail("登录失败已超过3次，系统被锁定，需要重启服务后才能继续")
	}

	config := configs.SharedAdminConfig()
	for _, user := range config.Users {
		if user.Username == params.Username && user.Password == params.Password {
			params.Auth.StoreUsername(user.Username)
			this.Next("/", nil, "").Success()
			return
		}
	}

	countLoginTries ++

	this.Fail("登录失败，请检查用户名密码")
}
