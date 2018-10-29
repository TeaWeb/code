package fastcgi

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
			Prefix("/proxy/fastcgi").
			Post("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			Post("/addParam", new(AddParamAction)).
			Post("/deleteParam", new(DeleteParamAction)).
			Post("/updateParam", new(UpdateParamAction)).
			Post("/updatePass", new(UpdatePassAction)).
			Post("/updateTimeout", new(UpdateTimeoutAction)).
			EndAll()
	})
}
