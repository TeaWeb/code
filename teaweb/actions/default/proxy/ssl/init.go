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
			GetPost("/update", new(UpdateAction)).
			Post("/startHttps", new(StartHttpsAction)).
			Post("/shutdownHttps", new(ShutdownHttpsAction)).
			Get("/downloadFile", new(DownloadFileAction)).
			Get("/generate", new(GenerateAction)).
			Get("/acmeCreateTask", new(AcmeCreateTaskAction)).
			GetPost("/acmeCreateUser", new(AcmeCreateUserAction)).
			Get("/acmeUsers", new(AcmeUsersAction)).
			Post("/acmeUserDelete", new(AcmeUserDeleteAction)).
			Post("/acmeRecords", new(AcmeRecordsAction)).
			Post("/acmeDnsChecking", new(AcmeDnsCheckingAction)).
			Post("/acmeDeleteTask", new(AcmeDeleteTaskAction)).
			Post("/acmeRenewTask", new(AcmeRenewTaskAction)).
			EndAll()
	})
}
