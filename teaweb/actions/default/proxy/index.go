package proxy

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
}) {
	servers := []maps.Map{}

	for _, config := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
		servers = append(servers, maps.Map{
			"config":   config,
			"filename": config.Filename,
		})
	}

	this.Data["servers"] = servers

	this.Show()
}
