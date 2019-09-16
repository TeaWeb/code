package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teadb"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

// 将Receiver转换为Map
func ConvertReceiversToMaps(receivers []*notices.NoticeReceiver) (result []maps.Map) {
	result = []maps.Map{}
	for _, receiver := range receivers {
		m := maps.Map{
			"name":      receiver.Name,
			"id":        receiver.Id,
			"user":      receiver.User,
			"mediaType": "",
		}

		// 媒介
		media := notices.SharedNoticeSetting().FindMedia(receiver.MediaId)
		if media != nil {
			m["mediaType"] = media.Name
		}
		result = append(result, m)
	}

	return result
}

// 发送一个后端下线通知
func NotifyProxyBackendDownMessage(serverId string, backend *teaconfigs.BackendConfig, location *teaconfigs.LocationConfig, websocket *teaconfigs.WebsocketConfig) error {
	level := notices.NoticeLevelWarning

	cond := notices.ProxyCond{
		ServerId:  serverId,
		BackendId: backend.Id,
		Level:     level,
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

	NotifyServer(serverId, level, "代理服务通知", "后端服务器'"+backend.Address+"'因错误过多已经下线")

	return nil
}

// 推送代理服务相关通知
func NotifyServer(serverId string, level notices.NoticeLevel, subject string, message string) {
	server := teaconfigs.NewServerConfigFromId(serverId)
	if server == nil {
		return
	}
	receivers := server.FindAllNoticeReceivers(level)
	if len(receivers) == 0 {
		return
	}
	noticeutils.AddTask(level, receivers, subject, message+"\n位置："+server.Description)
}
