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
			coll.CreateIndex(teamongo.NewIndexField("proxy.serverId", true))
			coll.CreateIndex(teamongo.NewIndexField("agent.agentId", true))
			coll.CreateIndex(
				teamongo.NewIndexField("agent.agentId", true),
				teamongo.NewIndexField("agent.appId", true),
				teamongo.NewIndexField("agent.itemId", true),
			)
		}()
	})
}
