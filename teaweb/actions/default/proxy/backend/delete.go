package backend

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename string
	Index    int
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(server.Backends) {
		server.Backends = lists.Remove(server.Backends, params.Index).([]*teaconfigs.ServerBackendConfig)
	}

	server.WriteToFilename(params.Filename)
	global.NotifyChange()

	this.Refresh().Success()
}
