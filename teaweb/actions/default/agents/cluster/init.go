package cluster

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Helper(new(Helper)).
			Prefix("/agents/cluster").
			Get("/add", new(AddAction)).
			Post("/search", new(SearchAction)).
			Post("/conn", new(ConnAction)).
			Post("/auth", new(AuthAction)).
			Post("/install", new(InstallAction)).
			EndAll()
	})
}
