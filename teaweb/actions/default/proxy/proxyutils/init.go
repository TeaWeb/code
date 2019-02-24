package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 启动统计
		go func() {
			for _, server := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
				if !server.On {
					continue
				}
				if server.StatBoard == nil && server.RealtimeBoard == nil {
					continue
				}
				ReloadServerStats(server.Id)
			}
		}()
	})
}
