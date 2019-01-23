package notices

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/agents/notices").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			Post("/setRead", new(SetReadAction)).
			EndAll()
	})
}
