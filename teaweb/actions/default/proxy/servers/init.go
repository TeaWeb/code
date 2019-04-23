package servers

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/servers").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Get("/waf", new(WafAction)).
			Post("/waf/update", new(WafUpdateAction)).
			EndAll()
	})
}
