package server

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(settings.Helper)).
			Prefix("/settings/server").
			Get("/http", new(HttpAction)).
			Post("/httpUpdate", new(HttpUpdateAction)).
			Get("/https", new(HttpsAction)).
			Post("/httpsUpdate", new(HttpsUpdateAction)).
			EndAll()
	})
}
