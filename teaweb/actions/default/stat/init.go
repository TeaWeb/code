package stat

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantStatistics,
			}).
			Helper(new(Helper)).
			Prefix("/stat").
			Get("", new(IndexAction)).
			Get("/data", new(DataAction)).
			EndAll()
	})
}
