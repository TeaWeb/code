package mongo

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/files"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	config := configs.SharedMongoConfig()

	this.Data["config"] = config
	this.Data["uri"] = config.URI()

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

	this.Show()
}
