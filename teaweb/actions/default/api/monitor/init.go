package monitor

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/api/monitor").
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
