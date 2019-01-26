package agent

import (
	"context"
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"io/ioutil"
	"sync"
	"time"
)

type PushAction actions.Action

// 接收推送的数据
func (this *PushAction) Run(params struct{}) {
	agent := this.Context.Get("agent").(*agents.AgentConfig)

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
			Id:         objectid.New(),
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
		appId := m.GetString("appId")
		itemId := m.GetString("itemId")
		app := agentutils.FindAgentApp(agent, appId)
		if app == nil {
			this.Success()
		}

		item := app.FindItem(itemId)
		if item == nil {
			this.Success()
		}

		v := m.Get("value")
		level, message := item.TestValue(v)

		// 通知消息
		if level != agents.NoticeLevelNone {
			notice := notices.NewNotice()
			notice.SetTime(t)
			notice.Message = message
			notice.Agent = notices.AgentCond{
				AgentId: agent.Id,
				AppId:   appId,
				ItemId:  itemId,
				Level:   level,
			}
			err = noticeutils.NewNoticeQuery().Insert(notice)
			if err != nil {
				logs.Error(err)
			}
		}

		value := &agents.Value{
			Id:          objectid.New(),
			AppId:       appId,
			AgentId:     agent.Id,
			ItemId:      itemId,
			Value:       v,
			Error:       m.GetString("error"),
			NoticeLevel: level,
		}
		value.SetTime(t)

		err := teamongo.NewValueQuery().Insert(value)
		if err != nil {
			logs.Error(err)
		}
	} else if eventDomain == "SystemAppsEvent" {
		result := struct {
			Apps []*agents.AppConfig
		}{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			logs.Error(err)
		} else {
			agentRuntime := agentutils.FindAgentRuntime(agent)
			agentRuntime.ResetSystemApps()
			agentRuntime.AddApps(result.Apps)
		}
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
