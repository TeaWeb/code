package tunnel

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/tunnel").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/generateSecret", new(GenerateSecretAction)).
			EndAll()
	})
}
