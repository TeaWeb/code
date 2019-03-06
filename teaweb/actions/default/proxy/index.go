package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type IndexAction actions.Action

// 代理首页
func (this *IndexAction) Run(params struct {
}) {
	// 跳转到第一个
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		logs.Error(err)
		return
	}
	servers := serverList.FindAllServers()
	if len(servers) > 0 {
		this.RedirectURL("/proxy/board?serverId=" + servers[0].Id)
	} else {
		this.Show()
	}
}
