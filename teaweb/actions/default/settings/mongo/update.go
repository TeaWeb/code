package mongo

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改连接
func (this *UpdateAction) Run(params struct{}) {
	config := configs.SharedMongoConfig()
	this.Data["config"] = configs.MongoConnectionConfig{
		Scheme:     config.Scheme,
		Username:   config.Username,
		Password:   "",
		Host:       config.Host,
		Port:       config.Port,
		RequestURI: config.RequestURI,
	}

	this.Show()
}

// 提交保存
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
	if len(params.Username) > 0 {
		if len(params.Password) > 0 {
			config.Password = params.Password
		}
	} else {
		config.Password = ""
	}
	err := config.Save()

	if err != nil {
		this.Fail("文件写入失败，请检查'configs/mongo.conf'写入权限")
	}

	// 重新连接
	teamongo.RestartClient()

	this.Next("/settings/mongo", nil).Success("保存成功")
}
