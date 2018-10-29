package monitor

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantApp,
			}).
			Helper(new(Helper)).
			Prefix("/monitor").
			Get("", new(IndexAction)).
			EndAll()
	})
}
