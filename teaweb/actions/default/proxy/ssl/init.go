package ssl

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
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
