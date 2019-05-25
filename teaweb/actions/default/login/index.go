package login

import (
	"github.com/TeaWeb/code/teaconfigs/audits"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"time"
)

type IndexAction actions.Action

// 登录
func (this *IndexAction) RunGet() {
	// 检查IP限制
	if !configs.SharedAdminConfig().AllowIP(this.RequestRemoteIP()) {
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		this.WriteString("TeaWeb Access Forbidden")
		return
	}

	b := Notify(this)
	if !b {
		return
	}

	this.Data["teaDemoEnabled"] = teaconst.DemoEnabled

	this.Show()
}

// 提交登录
func (this *IndexAction) RunPost(params struct {
	Username string
	Password string
	Remember bool
	Must     *actions.Must
	Auth     *helpers.UserShouldAuth
}) {
	// 记录
	teamongo.NewAuditsQuery().Insert(audits.NewLog(params.Username, audits.ActionLogin, "登录", map[string]string{
		"ip": this.RequestRemoteIP(),
	}))

	// 检查IP限制
	if !configs.SharedAdminConfig().AllowIP(this.RequestRemoteIP()) {
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		this.WriteString("TeaWeb Access Forbidden")
		return
	}

	b := Notify(this)
	if !b {
		return
	}

	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	adminConfig := configs.SharedAdminConfig()
	user := adminConfig.FindActiveUser(params.Username)
	if user != nil {
		// 错误次数
		if user.CountLoginTries() >= 3 {
			this.Fail("登录失败已超过3次，系统被锁定，需要重置服务后才能继续")
		}

		// 密码错误
		if user.Password != params.Password {
			user.IncreaseLoginTries()
			this.Fail("登录失败，请检查用户名密码")
		}

		user.ResetLoginTries()

		// Session
		params.Auth.StoreUsername(user.Username, params.Remember)

		// 记录登录IP
		user.LoggedAt = time.Now().Unix()
		user.LoggedIP = this.RequestRemoteIP()

		// 在开发环境下不保存登录IP，以便于不干扰git
		if !Tea.IsTesting() {
			adminConfig.Save()
		}

		this.Next("/", nil, "").Success()
		return
	}

	this.Fail("登录失败，请检查用户名密码")
}
