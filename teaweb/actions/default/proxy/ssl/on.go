package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type OnAction actions.Action

func (this *OnAction) Run(params struct {
	Filename string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if server.SSL == nil {
		ssl := new(teaconfigs.SSLConfig)
		ssl.On = true
		server.SSL = ssl
	} else {
		server.SSL.On = true
	}

	server.Save()

	proxyutils.NotifyChange()

	this.Success()
}
