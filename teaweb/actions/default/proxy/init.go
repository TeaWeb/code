package proxy

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Module("").
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/proxy").

			Get("", new(IndexAction)).
			Get("/status", new(StatusAction)).
			GetPost("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/detail", new(DetailAction)).
			Get("/httpOn", new(HttpOnAction)).
			Get("/httpOff", new(HttpOffAction)).
			Post("/updateDescription", new(UpdateDescriptionAction)).
			Post("/addName", new(AddNameAction)).
			Post("/updateName", new(UpdateNameAction)).
			Post("/deleteName", new(DeleteNameAction)).

			Post("/addListen", new(AddListenAction)).
			Post("/deleteListen", new(DeleteListenAction)).
			Post("/updateListen", new(UpdateListenAction)).

			Post("/updateRoot", new(UpdateRootAction)).

			Get("/restart", new(RestartAction)).

			EndAll()
	})
}
