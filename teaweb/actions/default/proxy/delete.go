package proxy

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) Run(params struct {
	Server string
}) {
	this.Data["server"] = maps.Map{
		"filename": params.Server,
	}

	this.Show()
}

func (this *DeleteAction) RunPost(params struct {
	Server string
}) {
	configFile := files.NewFile(Tea.ConfigFile(params.Server))
	if !configFile.Exists() {
		this.Fail("要删除的配置文件不存在")
	}

	err := configFile.Delete()
	if err != nil {
		logs.Error(err)
		this.Fail("配置文件删除失败")
	}

	// @TODO 删除对应的certificate file和certificate key file

	// 重启
	proxyutils.NotifyChange()

	this.Success()
}
