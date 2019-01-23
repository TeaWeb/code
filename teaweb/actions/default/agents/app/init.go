package app

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
			Prefix("/apps/app").
			Get("", new(IndexAction)).
			Get("/reload", new(ReloadAction)).
			Post("/favor", new(FavorAction)).
			Post("/cancelFavor", new(CancelFavorAction)).
			EndAll()
	})
}
