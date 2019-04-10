package mongo

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"runtime"
	"strings"
)

type IndexAction actions.Action

// MongoDB连接信息
func (this *IndexAction) Run(params struct{}) {
	config := configs.SharedMongoConfig()

	this.Data["config"] = maps.Map{
		"scheme":                  config.Scheme,
		"username":                config.Username,
		"password":                strings.Repeat("*", len(config.Password)),
		"host":                    config.Host,
		"port":                    config.Port,
		"authMechanism":           config.AuthMechanism,
		"authMechanismProperties": config.AuthMechanismPropertiesString(),
		"requestURI":              config.RequestURI,
	}
	this.Data["uri"] = config.URIMask()

	// 连接状态
	err := teamongo.Test()
	if err != nil {
		this.Data["error"] = err.Error()
	} else {
		this.Data["error"] = ""
	}

	// 检测是否已安装
	mongodbPath := Tea.Root + "/mongodb/bin/mongod"
	if files.NewFile(mongodbPath).Exists() {
		this.Data["isInstalled"] = true
	} else {
		this.Data["isInstalled"] = false
	}

	// 当前系统
	this.Data["os"] = runtime.GOOS

	this.Show()
}
