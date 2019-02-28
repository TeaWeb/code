package agentutils

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"strings"
	"time"
)

func init() {
	// 检查Agent连通性
	checkConnecting()
}

// 检查Agent连通性
func checkConnecting() {
	timers.Loop(60*time.Second, func(looper *timers.Looper) {
		agentList, err := agents.SharedAgentList()
		if err != nil {
			return
		}
		for _, agent := range agentList.FindAllAgents() {
			if !agent.On {
				continue
			}

			runtimeAgent := FindAgentRuntime(agent)

			// 监控连通性
			isWaiting := CheckAgentIsWaiting(agent.Id)
			if !isWaiting {
				runtimeAgent.CountDisconnections ++

				if runtimeAgent.CountDisconnections >= 3 { // 失去连接3次则提醒
					runtimeAgent.CountDisconnections = 0
					sendDisconnectNotice(runtimeAgent)
				}
			} else {
				runtimeAgent.CountDisconnections = 0
			}
		}
	})
}

// 发送Agent失联通知
func sendDisconnectNotice(agent *agents.AgentConfig) {
	message := "Agent\"" + agent.Name + "\"失去连接"
	level := notices.NoticeLevelError
	t := time.Now()

	notice := notices.NewNotice()
	notice.Id = objectid.New()
	notice.Agent.AgentId = agent.Id
	notice.Agent.Level = level
	notice.Message = message
	notice.SetTime(t)
	err := noticeutils.NewNoticeQuery().Insert(notice)
	if err != nil {
		logs.Error(err)
	} else {
		// 通过媒介发送通知
		setting := notices.SharedNoticeSetting()
		fullMessage := "消息：" + message + "\n时间：" + timeutil.Format("Y-m-d H:i:s", t)
		linkNames := []string{}
		for _, l := range FindNoticeLinks(notice) {
			linkNames = append(linkNames, types.String(l["name"]))
		}
		if len(linkNames) > 0 {
			fullMessage += "\n位置：" + strings.Join(linkNames, "/")
		}
		receiverIds := setting.Notify(level, fullMessage, func(receiverId string, minutes int) int {
			return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
				"agent.agentId": agent.Id,
				"agent.appId":   "",
			}, minutes)
		})
		if len(receiverIds) > 0 {
			noticeutils.UpdateNoticeReceivers(notice.Id, receiverIds)
		}
	}
}
