package mongo

import (
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/TeaWeb/code/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改连接
func (this *UpdateAction) Run(params struct{}) {
	config, err := db.LoadMongoConfig()
	if err != nil {
		config = db.DefaultMongoConfig()
	}
	this.Data["config"] = maps.Map{
		"scheme":                  config.Scheme,
		"username":                config.Username,
		"password":                "",
		"host":                    config.Host(),
		"port":                    config.Port(),
		"dbName":                  config.DBName,
		"authEnabled":             config.AuthEnabled,
		"authMechanism":           config.AuthMechanism,
		"authMechanismProperties": config.AuthMechanismPropertiesString(),
		"requestURI":              config.RequestURI,
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	Host                    string
	Port                    uint
	DBName                  string `alias:"dbName"`
	Username                string
	Password                string
	AuthEnabled             bool
	AuthMechanism           string
	AuthMechanismProperties string

	Must *actions.Must
}) {
	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Require("请输入端口").
		Gt(0, "请输入正确的端口")

	config, err := db.LoadMongoConfig()
	if err != nil {
		this.Fail(err.Error())
	}

	config.SetAddr(params.Host, params.Port)
	config.DBName = params.DBName
	config.AuthEnabled = params.AuthEnabled
	config.AuthMechanism = params.AuthMechanism
	config.LoadAuthMechanismProperties(params.AuthMechanismProperties)
	config.Username = params.Username
	if len(params.Username) > 0 {
		if len(params.Password) > 0 {
			config.Password = params.Password
		}
	} else {
		config.Password = ""
	}
	err = config.Save()

	if err != nil {
		this.Fail("文件写入失败，请检查'configs/mongo.conf'写入权限")
	}

	// 重新连接
	teadb.ChangeDB()

	this.Next("/settings/mongo", nil).Success("保存成功")
}
