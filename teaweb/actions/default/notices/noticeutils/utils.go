package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teadb"
	"github.com/iwind/TeaGo/logs"
)

// 发送一个后端下线通知
func NotifyProxyBackendDownMessage(serverId string, backend *teaconfigs.BackendConfig, location *teaconfigs.LocationConfig, websocket *teaconfigs.WebsocketConfig) error {
	cond := notices.ProxyCond{
		ServerId:  serverId,
		BackendId: backend.Id,
		Level:     notices.NoticeLevelWarning,
	}
	if location != nil {
		cond.LocationId = location.Id
	}
	if websocket != nil {
		cond.Websocket = true
	}

	// 不阻塞
	go func() {
		err := teadb.NoticeDAO().NotifyProxyMessage(cond, "后端服务器'"+backend.Address+"'因错误过多已经下线")
		if err != nil {
			logs.Error(err)
		}
	}()
	return nil
}
