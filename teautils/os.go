package teautils

import (
	"github.com/iwind/TeaGo/Tea"
)

// 临时文件
func TmpFile(path string) string {
	return Tea.Root + Tea.DS + "web" + Tea.DS + "tmp" + Tea.DS + path
}
