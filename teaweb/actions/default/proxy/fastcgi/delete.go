package fastcgi

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除Fastcgi
func (this *DeleteAction) Run(params struct {
	Server     string
	LocationId string
	FastcgiId  string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgiList.RemoveFastcgi(params.FastcgiId)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
