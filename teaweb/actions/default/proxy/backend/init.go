package backend

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
			Module("").
			Prefix("/proxy/backend").
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/scheduling", new(SchedulingAction)).
			Post("/online", new(OnlineAction)).
			Post("/clearFails", new(ClearFailsAction)).
			Prefix("").
			EndAll()
	})
}
