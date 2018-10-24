package apps

import (
	_ "github.com/TeaWeb/code/teaweb/actions/default/apps/app"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/apps").
			Get("", new(IndexAction)).
			Get("/all", new(AllAction)).
			EndAll()
	})
}
