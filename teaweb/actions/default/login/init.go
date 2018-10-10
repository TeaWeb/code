package login

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			EndAll().
			Prefix("/login").
			GetPost("", new(IndexAction)).
			Prefix("").
			EndAll()
	})
}
