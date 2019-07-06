// +build !windows

package teaweb

// 启动服务模式
func (this *WebShell) execService() bool {
	// do nothing beyond windows
	return true
}
