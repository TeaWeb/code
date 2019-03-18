package noticeutils

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 为notices创建索引
		go func() {
			coll := teamongo.FindCollection("notices")
			coll.CreateIndex(map[string]bool{
				"proxy.serverId": true,
			})
			coll.CreateIndex(map[string]bool{
				"agent.agentId": true,
			})
			coll.CreateIndex(map[string]bool{
				"agent.agentId": true,
				"agent.appId":   true,
				"agent.itemId":  true,
			})
		}()
	})
}
