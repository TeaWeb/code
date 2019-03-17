package cluster

import "github.com/iwind/TeaGo/actions"

type AuthAction actions.Action

// 认证
func (this *AuthAction) Run(params struct {
	Master   string
	Dir      string
	Username string
	Password string
	Must     *actions.Must
}) {
	params.Must.
		Field("master", params.Master).
		Require("请输入TeaWeb访问地址").
		Field("dir", params.Dir).
		Require("请输入安装目录").
		Field("username", params.Username).
		Require("请输入登录主机的用户名").
		Field("password", params.Password).
		Require("请输入登录主机的密码")

	this.Success()
}
