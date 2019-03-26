package apps

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Helper(new(Helper)).
			Prefix("/agents/board").
			GetPost("", new(IndexAction)).
			Get("/charts", new(ChartsAction)).
			Post("/addChart", new(AddChartAction)).
			Post("/removeChart", new(RemoveChartAction)).
			Post("/moveChart", new(MoveChartAction)).
			Post("/initDefaultApp", new(InitDefaultAppAction)).
			EndAll()
	})
}
