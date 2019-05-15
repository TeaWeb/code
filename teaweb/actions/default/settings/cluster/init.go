package cluster

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/settings/cluster").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/connect", new(ConnectAction)).
			Post("/sync", new(SyncAction)).
			EndAll()

	})
}
