package ssl

import (
	"github.com/TeaWeb/code/teacluster"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/ssl/sslutils"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

func init() {
	// 路由定义
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
			Get("/acmeTask", new(AcmeTaskAction)).
			Get("/acmeDownload", new(AcmeDownloadAction)).
			EndAll()
	})

	// 检查ACME证书更新
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		timers.Every(24*time.Hour, func(ticker *time.Ticker) {
			renewACMECerts()
		})
	})
}

// 检查ACME证书更新
func renewACMECerts() {
	logs.Println("[acme]check acme requests")

	// skip slave node
	node := teaconfigs.SharedNodeConfig()
	if node != nil && node.On && !node.IsMaster() {
		return
	}
	nodeDirty := false
	if node != nil && teacluster.SharedManager.IsActive() {
		// 集群节点状态
		nodeDirty = teacluster.SharedManager.IsChanged()
	}

	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		return
	}

	nodeIsChanged := false

	for _, server := range serverList.FindAllServers() {
		if server.SSL == nil || !server.SSL.On || len(server.SSL.CertTasks) == 0 {
			continue
		}
		serverIsChanged := false
		for _, task := range server.SSL.CertTasks {
			if !task.On {
				continue
			}
			if task.Request == nil {
				continue
			}
			date := task.Request.CertDate()
			if len(date[1]) == 0 {
				continue
			}
			if timeutil.Format("Y-m-d") >= date[1] {
				client, err := task.Request.Client()
				if err != nil {
					task.RunAt = time.Now().Unix()
					task.RunError = err.Error()
					logs.Error(err)
					serverIsChanged = true
					continue
				}
				err = task.Request.Renew(client)
				if err != nil {
					task.RunAt = time.Now().Unix()
					task.RunError = err.Error()
					logs.Error(err)
					serverIsChanged = true
					continue
				}

				task.RunAt = time.Now().Unix()
				task.RunError = ""
				serverIsChanged = true

				// 更新证书
				found := false
				for _, cert := range server.SSL.Certs {
					if cert.TaskId == task.Id {
						err = task.Request.WriteCertFile(Tea.ConfigFile(cert.CertFile))
						if err != nil {
							logs.Error(err)
						}

						err = task.Request.WriteKeyFile(Tea.ConfigFile(cert.KeyFile))
						if err != nil {
							logs.Error(err)
						}

						found = true
					}
				}

				// 重新加载证书
				if found {
					errs := sslutils.ReloadACMECert(server.Id, task.Id)
					for _, err2 := range errs {
						logs.Println("[acme]reload acme task:", err2.Error())
					}
				}
			}
		}

		if serverIsChanged {
			nodeIsChanged = true
			server.Save()
		}
	}

	// 如果先前节点没有变更，则自动推送到集群
	if !nodeDirty && nodeIsChanged {
		node := teaconfigs.SharedNodeConfig()
		if node != nil && node.On && node.IsMaster() && teacluster.SharedManager.IsActive() {
			teacluster.SharedManager.PushItems()
			teacluster.SharedManager.SetIsChanged(false)
		}
	}
}
