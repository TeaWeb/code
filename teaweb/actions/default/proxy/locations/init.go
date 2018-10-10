package locations

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/proxy/locations").
			Helper(new(helpers.UserMustAuth)).
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
			EndAll()
	})
}
