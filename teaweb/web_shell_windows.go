// +build windows

package teaweb

import "github.com/TeaWeb/code/teautils"

// 启动服务模式
func (this *WebShell) execService() bool {
	// start the manager
	manager := teautils.NewServiceManager("TeaWeb", "TeaWeb Server")
	manager.Run()

	return true
}
