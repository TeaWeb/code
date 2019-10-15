package database

import (
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Helper(new(settings.Helper)).
			Prefix("/settings/database").
			Get("", new(IndexAction)).
			Post("/tables", new(TablesAction)).
			Post("/tableStat", new(TableStatAction)).
			Post("/deleteTable", new(DeleteTableAction)).
			EndAll()
	})
}
