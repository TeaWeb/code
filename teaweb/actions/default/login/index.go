package login

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/audits"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teadb"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
	"time"
)

type IndexAction actions.Action

var TokenSalt = stringutil.Rand(32)

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

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	this.Data["token"] = stringutil.Md5(TokenSalt+timestamp) + timestamp

	this.Show()
}

// 提交登录
func (this *IndexAction) RunPost(params struct {
	Username string
	Password string
	Token    string
	Remember bool
	Must     *actions.Must
	Auth     *helpers.UserShouldAuth
}) {
	// 记录登录
	go func() {
		err := teadb.AuditLogDAO().InsertOne(audits.NewLog(params.Username, audits.ActionLogin, "登录", map[string]string{
			"ip": this.RequestRemoteIP(),
		}))
		if err != nil {
			logs.Error(err)
		}
	}()

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

	if params.Password == stringutil.Md5("") {
		this.FailField("password", "请输入密码")
	}

	// 检查token
	if len(params.Token) <= 32 {
		this.Fail("请通过登录页面登录")
	}
	timestampString := params.Token[32:]
	if stringutil.Md5(TokenSalt+timestampString) != params.Token[:32] {
		this.FailField("refresh", "登录页面已过期，请刷新后重试")
	}
	timestamp := types.Int64(timestampString)
	if timestamp < time.Now().Unix()-1800 {
		this.FailField("refresh", "登录页面已过期，请刷新后重试")
	}

	// 查找用户
	adminConfig := configs.SharedAdminConfig()
	user := adminConfig.FindActiveUser(params.Username)
	if user != nil {
		// 错误次数
		if user.CountLoginTries() >= 3 {
			this.Fail("登录失败已超过3次，系统被锁定，需要重置服务后才能继续")
		}

		// 密码错误
		if !adminConfig.ComparePassword(params.Password, user.Password) {
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
			err := adminConfig.Save()
			if err != nil {
				logs.Error(err)
			}
		}

		this.Next("/", nil, "").Success()
		return
	}

	this.Fail("登录失败，请检查用户名密码")
}
