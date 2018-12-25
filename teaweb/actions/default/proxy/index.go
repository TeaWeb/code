package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 代理首页
func (this *IndexAction) Run(params struct {
}) {
	// 跳转到第一个
	servers := teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir())
	if len(servers) > 0 {
		this.RedirectURL("/proxy/board?server=" + servers[0].Filename)
	} else {
		this.Show()
	}
}
