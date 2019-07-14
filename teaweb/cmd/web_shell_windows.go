// +build windows

package cmd

import (
	"github.com/TeaWeb/code/teautils"
	"io"
)

// 启动服务模式
func (this *WebShell) ExecService(writer io.Writer) bool {
	// start the manager
	manager := teautils.NewServiceManager("TeaWeb", "TeaWeb Server")
	manager.Run()

	return true
}
