package fastcgi

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Prefix("/proxy/fastcgi").
			Get("/data", new(DataAction)).
			GetPost("/add", new(AddAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			EndAll()
	})
}
