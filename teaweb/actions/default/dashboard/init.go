package dashboard

import (
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Prefix("/dashboard").
			GetPost("", new(IndexAction)).
			Get("/logs", new(LogsAction)).
			EndAll()
	})
}
