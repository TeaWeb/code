package backend

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Module("").
			Prefix("/proxy/backend").
			Get("", new(IndexAction)).
			Post("/add", new(AddAction)).
			Post("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Prefix("").
			EndAll()
	})
}
