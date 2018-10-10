package log

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			EndAll().
			Helper(new(helpers.UserMustAuth)).
			Prefix("/log").
			Get("", new(IndexAction)).
			Get("/get", new(GetAction)).
			EndAll()
	})
}
