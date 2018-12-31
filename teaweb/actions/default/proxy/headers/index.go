package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 自定义Http Header
func (this *IndexAction) Run(params struct {
	Server string // 必填
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "header"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = proxy

	this.Data["headerQuery"] = maps.Map{
		"server": params.Server,
	}

	this.Show()
}
