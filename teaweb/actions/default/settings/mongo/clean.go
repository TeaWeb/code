package mongo

import (
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/iwind/TeaGo/actions"
)

type CleanAction actions.Action

// 设置自动清理
func (this *CleanAction) Run(params struct{}) {
	config, _ := db.LoadMongoConfig()
	this.Data["accessLog"] = config.AccessLog

	this.Show()
}
