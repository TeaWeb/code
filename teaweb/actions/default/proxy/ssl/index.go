package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// SSL设置
func (this *IndexAction) Run(params struct {
	Server string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "https"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = proxy

	this.Show()
}
