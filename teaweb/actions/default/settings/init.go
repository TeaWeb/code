package settings

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
)

var serverChanged = false

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/settings").
			Get("", new(IndexAction)).
			EndAll()
	})
}

func NotifyServerChange() {
	serverChanged = true
}

func ServerIsChanged() bool {
	return serverChanged
}
