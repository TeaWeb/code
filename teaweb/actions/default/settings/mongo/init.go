package mongo

import (
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Helper(new(settings.Helper)).
			Prefix("/settings/mongo").
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/test", new(TestAction)).
			GetPost("/install", new(InstallAction)).
			Get("/installStatus", new(InstallStatusAction)).
			EndAll()
	})
}
