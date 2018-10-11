package mongo

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teamongo"
)

type UpdateAction actions.Action

func (this *UpdateAction) Run(params struct{}) {
	this.Data["config"] = configs.SharedMongoConfig()

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	Host     string
	Port     uint
	Username string
	Password string

	Must *actions.Must
}) {
	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Require("请输入端口").
		Gt(0, "请输入正确的端口")

	config := configs.SharedMongoConfig()
	config.Host = params.Host
	config.Port = params.Port
	config.Username = params.Username
	config.Password = params.Password
	err := config.WriteBack()

	if err != nil {
		this.Fail("文件写入失败，请检查'configs/mongo.conf'写入权限")
	}

	// 重新连接
	teamongo.RestartClient()

	this.Next("/settings/mongo", nil).Success("保存成功")
}
