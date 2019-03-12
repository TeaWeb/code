package agent

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"time"
)

type PullAction actions.Action

// 拉取事件
func (this *PullAction) Run(params struct{}) {
	agentId := this.Context.Get("agent").(*agents.AgentConfig).Id
	agentVersion := this.Request.Header.Get("Tea-Agent-Version")
	nano := this.Request.Header.Get("Tea-Agent-Nano")
	speed := float64(0)
	if len(nano) > 0 {
		speed = float64(time.Now().UnixNano()-types.Int64(nano)) / 1000000
	}

	c := make(chan *agentutils.Event)
	agentutils.WaitAgentQueue(agentId, agentVersion, speed, c)

	// 监控是否中断请求
	go func() {
		<-this.Request.Context().Done()
		agentutils.RemoveAgentQueue(agentId, c)

		// 关闭channel
		c <- nil
		close(c)
	}()

	events := []*agentutils.Event{}
	for {
		event := <-c
		func() {
			if event != nil {
				events = append(events, event)
			}
		}()

		break
	}

	// 移除
	agentutils.RemoveAgentQueue(agentId, c)

	this.Data["events"] = events

	this.Success()
}
