package sslutils

import (
	"github.com/TeaWeb/code/teaproxy"
)

// 重载证书
func ReloadACMECert(serverId string, taskId string) (errs []error) {
	if len(serverId) == 0 || len(taskId) == 0 {
		return
	}

	// 更新目前正在使用的Server
	runningServer := teaproxy.SharedManager.FindServer(serverId)
	if runningServer == nil || runningServer.SSL == nil || len(runningServer.SSL.Certs) == 0 {
		return
	}

	// 查找正在使用此任务的证书
	for _, cert := range runningServer.SSL.Certs {
		if cert.TaskId == taskId {
			err := cert.Validate()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}
