package rewrite

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/rewrite").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Post("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Post("/update", new(UpdateAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			EndAll()
	})
}
