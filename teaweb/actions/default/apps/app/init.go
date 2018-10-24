package app

import (
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Prefix("/apps/app").
			Get("", new(IndexAction)).
			Get("/reload", new(ReloadAction)).
			EndAll()
	})
}
