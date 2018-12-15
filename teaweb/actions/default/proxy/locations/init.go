package locations

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/headers"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/proxy/locations").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Get("", new(IndexAction)).
			Post("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Post("/moveUp", new(MoveUpAction)).
			Post("/moveDown", new(MoveDownAction)).
			Get("/detail", new(DetailAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			Post("/updateReverse", new(UpdateReverseAction)).
			Post("/updateCaseInsensitive", new(UpdateCaseInsensitiveAction)).
			Post("/updatePattern", new(UpdatePatternAction)).
			Post("/updateType", new(UpdateTypeAction)).
			Post("/updateRoot", new(UpdateRootAction)).
			Post("/updateCharset", new(UpdateCharsetAction)).
			Post("/updateIndex", new(UpdateIndexAction)).
			Post("/updateCache", new(UpdateCacheAction)).
			EndAll()
	})
}
