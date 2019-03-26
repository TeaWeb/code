package mongo

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type CleanAction actions.Action

// 设置自动清理
func (this *CleanAction) Run(params struct{}) {
	config, _ := configs.LoadMongoConfig()
	this.Data["accessLog"] = config.AccessLog

	this.Show()
}
