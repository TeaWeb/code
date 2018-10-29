package ssl

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Module("").
			Prefix("/proxy/ssl").
			Get("", new(IndexAction)).
			Post("/uploadCert", new(UploadCertAction)).
			Post("/uploadKey", new(UploadKeyAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			Prefix("").
			EndAll()
	})
}
