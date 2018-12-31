package rewrite

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除重写规则
func (this *DeleteAction) Run(params struct {
	Server     string
	LocationId string
	RewriteId  string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}
	rewriteList, err := server.FindRewriteList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	rewriteList.RemoveRewriteRule(params.RewriteId)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
