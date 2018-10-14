package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"time"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
}) {
	servers := []maps.Map{}

	for _, config := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
		servers = append(servers, maps.Map{
			"config":   config,
			"filename": config.Filename,

			// 10分钟内是否有访问
			"isActive": tealogs.SharedLogger().CountSuccessLogs(time.Now().Add(-10 * time.Minute).Unix(), time.Now().Unix(), config.Id) > 0,
		})
	}

	this.Data["servers"] = servers

	this.Show()
}
