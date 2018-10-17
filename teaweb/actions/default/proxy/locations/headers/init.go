package headers

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(proxy.Helper)).
			Prefix("/proxy/locations/headers").
			Post("/set", new(SetAction)).
			Post("/delete", new(DeleteAction)).
			Post("/update", new(UpdateAction)).
			Post("/addIgnore", new(AddIgnoreAction)).
			Post("/updateIgnore", new(UpdateIgnoreAction)).
			Post("/deleteIgnore", new(DeleteIgnoreAction)).
			EndAll()
	})
}
