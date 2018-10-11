package mongo

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teamongo"
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

	this.Show()
}
