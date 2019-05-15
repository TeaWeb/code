package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type PushAction actions.Action

// 接收推送的数据
func (this *PushAction) Run(params struct{}) {
	agent := this.Context.Get("agent").(*agents.AgentConfig)

	// 是否未启用
	if !agent.On {
		this.Success()
	}

	data, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		logs.Error(err)
		this.Fail("read body error")
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		logs.Error(err)
		this.Fail("unmarshal error")
	}

	timestamp := m.GetInt64("timestamp")
	t := time.Unix(timestamp, 0)

	eventDomain := m.GetString("event")

	if eventDomain == "ProcessEvent" { // 进程事件
		event := agentutils.ProcessLog{
			Id:         primitive.NewObjectID(),
			AgentId:    agent.Id,
			TaskId:     m.GetString("taskId"),
			ProcessId:  m.GetString("uniqueId"),
			ProcessPid: m.GetInt("pid"),
			EventType:  m.GetString("eventType"),
			Data:       m.GetString("data"),
			Timestamp:  timestamp,
			TimeFormat: struct {
				Year   string `bson:"year" json:"year"`
				Month  string `bson:"month" json:"month"`
				Day    string `bson:"day" json:"day"`
				Hour   string `bson:"hour" json:"hour"`
				Minute string `bson:"minute" json:"minute"`
				Second string `bson:"second" json:"second"`
			}{
				Year:   timeutil.Format("Y", t),
				Month:  timeutil.Format("Ym", t),
				Day:    timeutil.Format("Ymd", t),
				Hour:   timeutil.Format("YmdH", t),
				Minute: timeutil.Format("YmdHi", t),
				Second: timeutil.Format("YmdHis", t),
			},
		}

		coll := this.selectProcessEventCollection(agent.Id)
		_, err = coll.InsertOne(context.Background(), event)
		if err != nil {
			logs.Error(err)
		}
	} else if eventDomain == "ItemEvent" { // 监控项事件
		this.processItemEvent(agent, m, t)
	}

	this.Success()
}

var agentCollectionMap = map[string]*teamongo.Collection{} // agentId => collection
var agentCollectionLocker = sync.Mutex{}

func (this *PushAction) selectProcessEventCollection(agentId string) *teamongo.Collection {
	createdNew := false

	agentCollectionLocker.Lock()
	coll, found := agentCollectionMap[agentId]
	if !found {
		createdNew = true

		coll = teamongo.FindCollection("logs.agent." + agentId)
		agentCollectionMap[agentId] = coll
	}
	agentCollectionLocker.Unlock()

	if createdNew {
		coll.CreateIndex(map[string]bool{
			"agentId": true,
		})
		coll.CreateIndex(map[string]bool{
			"taskId": true,
		})
	}

	return coll
}

// 处理监控项事件
func (this *PushAction) processItemEvent(agent *agents.AgentConfig, m maps.Map, t time.Time) {
	appId := m.GetString("appId")
	itemId := m.GetString("itemId")
	app := agent.FindApp(appId)
	if app == nil {
		this.Success()
	}

	item := app.FindItem(itemId)
	if item == nil {
		this.Success()
	}

	v := m.Get("value")
	oldValue, err := this.findLatestAgentValue(agent.Id, appId, itemId)
	if err != nil {
		if err != context.DeadlineExceeded {
			logs.Error(err)
		}
		return
	}
	if oldValue == nil {
		oldValue = v
	}
	threshold, row, level, message, err := item.TestValue(v, oldValue)
	if err != nil {
		logs.Error(errors.New(item.Name + " " + err.Error()))
		if len(m.GetString("error")) == 0 {
			m["error"] = err.Error()
		}
	}

	// 处理消息中的变量
	message = teaconfigs.RegexpNamedVariable.ReplaceAllStringFunc(message, func(s string) string {
		result, err := agents.EvalParam(s, v, oldValue, maps.Map{
			"AGENT": maps.Map{
				"name": agent.Name,
				"host": agent.Host,
			},
			"APP": maps.Map{
				"name": app.Name,
			},
			"ITEM": maps.Map{
				"name": item.Name,
			},
			"ROW": row,
		}, false)
		if err != nil {
			logs.Error(err)
		}
		return result
	})

	if threshold != nil && len(message) == 0 {
		message = threshold.Expression()
	}

	// 通知消息
	setting := notices.SharedNoticeSetting()

	isNotified := false

	if level != notices.NoticeLevelNone {
		// 是否发送通知
		shouldNotify := true

		// 检查最近N此数值是否都是同类错误
		if threshold != nil && threshold.MaxFails > 1 {
			query := teamongo.NewAgentValueQuery()
			query.Agent(agent.Id)
			query.App(app.Id)
			query.Item(item.Id)
			query.Desc("_id")
			query.Limit(int64(threshold.MaxFails) - 1)
			values, err := query.FindAll()
			if err != nil {
				logs.Error(err)
			} else {
				if len(values) != threshold.MaxFails-1 { // 未达到连续失败次数
					shouldNotify = false
				} else {
					for _, v := range values {
						if v.ThresholdId != threshold.Id || v.IsNotified {
							shouldNotify = false
							break
						}
					}
				}
			}
		}

		// 发送通知
		if shouldNotify {
			isNotified = true

			notice := notices.NewNotice()
			notice.SetTime(t)
			notice.Message = message
			notice.Agent = notices.AgentCond{
				AgentId: agent.Id,
				AppId:   appId,
				ItemId:  itemId,
				Level:   level,
			}
			if threshold != nil {
				notice.Agent.Threshold = threshold.Expression()
			}
			notice.Hash()

			if notices.IsFailureLevel(level) {
				// 同样的消息短时间内只发送一条
				if noticeutils.ExistNoticesWithHash(notice.MessageHash, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   appId,
					"agent.itemId":  itemId,
				}, 1*time.Hour) {
					shouldNotify = false
				}
			}

			if shouldNotify {
				err := noticeutils.NewNoticeQuery().Insert(notice)
				if err != nil {
					logs.Error(err)
				}

				// 通过媒介发送通知
				fullMessage := "消息：" + message + "\n时间：" + timeutil.Format("Y-m-d H:i:s", t)
				linkNames := []string{}
				for _, l := range agentutils.FindNoticeLinks(notice) {
					linkNames = append(linkNames, types.String(l["name"]))
				}
				if len(linkNames) > 0 {
					fullMessage += "\n位置：" + strings.Join(linkNames, "/")
				}

				receiverIds := this.notifyMessage(agent, appId, itemId, setting, level, "有新的通知", fullMessage, false)
				if len(receiverIds) > 0 {
					noticeutils.UpdateNoticeReceivers(notice.Id, receiverIds)
				}
			}
		}
	}

	// 数值记录
	node := teaconfigs.SharedNodeConfig()
	nodeId := ""
	if node != nil {
		nodeId = node.Id
	}
	value := &agents.Value{
		Id:          primitive.NewObjectID(),
		NodeId:      nodeId,
		AppId:       appId,
		AgentId:     agent.Id,
		ItemId:      itemId,
		Value:       v,
		Error:       m.GetString("error"),
		NoticeLevel: level,
		CreatedAt:   time.Now().Unix(),
		IsNotified:  isNotified,
	}
	if threshold != nil {
		value.ThresholdId = threshold.Id
		value.Threshold = threshold.Expression()
	}
	value.SetTime(t)

	err = teamongo.NewAgentValueQuery().Insert(value)
	if err != nil {
		logs.Error(err)
		return
	}

	// 是否发送恢复通知
	if len(item.Id) > 0 && !notices.IsFailureLevel(level) {
		recoverSuccesses := item.RecoverSuccesses
		if recoverSuccesses <= 0 {
			recoverSuccesses = 1
		}

		query := teamongo.NewAgentValueQuery()
		query.Agent(agent.Id)
		query.App(app.Id)
		query.Item(item.Id)
		query.Desc("_id")
		query.Limit(int64(recoverSuccesses + 1))
		values, err := query.FindAll()
		if err != nil {
			logs.Error(err)
			return
		}

		lists.Reverse(values)

		countValues := len(values)
		if countValues != recoverSuccesses+1 {
			return
		}

		if !notices.IsFailureLevel(values[0].NoticeLevel) {
			return
		}

		success := true
		for i := 1; i < countValues; i ++ {
			if notices.IsFailureLevel(values[i].NoticeLevel) {
				success = false
				break
			}
		}

		if success {
			// 发送成功级别的通知
			notice := notices.NewNotice()
			notice.SetTime(t)
			notice.Agent = notices.AgentCond{
				AgentId: agent.Id,
				AppId:   appId,
				ItemId:  itemId,
				Level:   notices.NoticeLevelSuccess,
			}
			notice.Message = "监控项经过" + fmt.Sprintf("%d", recoverSuccesses) + "次刷新后，判定已恢复正常"
			linkNames := []string{}
			for _, l := range agentutils.FindNoticeLinks(notice) {
				linkNames = append(linkNames, types.String(l["name"]))
			}
			if len(linkNames) > 0 {
				notice.Message += " \n位置：" + strings.Join(linkNames, "/")
			}
			notice.Hash()
			err := noticeutils.NewNoticeQuery().Insert(notice)
			if err != nil {
				logs.Error(err)
			}

			this.notifyMessage(agent, appId, itemId, setting, notices.NoticeLevelSuccess, "有新的通知", notice.Message, true)
		}
	}
}

// 通知消息
func (this *PushAction) notifyMessage(agent *agents.AgentConfig, appId string, itemId string, setting *notices.NoticeSetting, level notices.NoticeLevel, subject string, message string, isSuccess bool) []string {

	isNotified := false
	receiverIds := []string{}

	receiverLevels := []notices.NoticeLevel{level}
	if isSuccess {
		receiverLevels = append(receiverLevels, notices.NoticeLevelError, notices.NoticeLevelWarning)
	}

	// 查找App的通知设置
	app := agent.FindApp(appId)
	if app != nil {
		receivers := app.FindAllNoticeReceivers(receiverLevels...)
		if len(receivers) > 0 {
			isNotified = true
			receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]"+subject, message, func(receiverId string, minutes int) int {
				return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   appId,
					"agent.itemId":  itemId,
				}, minutes)
			})
		}
	}

	// 查找Agent的通知设置
	if !isNotified {
		receivers := agent.FindAllNoticeReceivers(receiverLevels ...)
		if len(receivers) > 0 {
			isNotified = true
			receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]"+subject, message, func(receiverId string, minutes int) int {
				return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   appId,
					"agent.itemId":  itemId,
				}, minutes)
			})
		}
	}

	// 查找分组的通知设置
	if !isNotified {
		groupId := ""
		if len(agent.GroupIds) > 0 {
			groupId = agent.GroupIds[0]
		}
		group := agents.SharedGroupConfig().FindGroup(groupId)
		if group != nil {
			receivers := group.FindAllNoticeReceivers(receiverLevels...)
			if len(receivers) > 0 {
				isNotified = true
				receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]"+subject, message, func(receiverId string, minutes int) int {
					return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
						"agent.agentId": agent.Id,
						"agent.appId":   appId,
						"agent.itemId":  itemId,
					}, minutes)
				})
			}
		}
	}

	// 全局通知
	if !isNotified {
		receivers := setting.FindAllNoticeReceivers(receiverLevels...)
		if len(receivers) > 0 {
			receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]"+subject, message, func(receiverId string, minutes int) int {
				return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   appId,
					"agent.itemId":  itemId,
				}, minutes)
			})
		}
	}

	return receiverIds
}

// 查找最近的一次数值记录
func (this *PushAction) findLatestAgentValue(agentId string, appId string, itemId string) (interface{}, error) {
	query := teamongo.NewAgentValueQuery()
	query.Agent(agentId)
	query.App(appId)
	query.Item(itemId)
	query.Attr("error", "")
	query.Desc("_id")
	v, err := query.Find()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.Value, nil
}
