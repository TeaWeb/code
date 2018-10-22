package teaplugins

import (
	"github.com/iwind/TeaGo/Tea"
)

func init() {
	if Tea.IsTesting() {
		return
	}

	load()
}
