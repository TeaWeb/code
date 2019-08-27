package agentutils

import (
	"sync"
)

// Agent事件
type AgentEvent struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

var eventQueueMap = map[string]map[chan *AgentEvent]*AgentState{} // agentId => { chan => State }
var eventQueueLocker = sync.Mutex{}

// 新Agent事件
func NewAgentEvent(name string, data interface{}) *AgentEvent {
	return &AgentEvent{
		Name: name,
		Data: data,
	}
}

// 等待Agent事件
func WaitAgentQueue(agentId string, agentVersion string, osName string, speed float64, ip string, c chan *AgentEvent) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	_, found := eventQueueMap[agentId]
	if found {
		eventQueueMap[agentId][c] = &AgentState{
			Version:     agentVersion,
			OsName:      osName,
			Speed:       speed,
			IP:          ip,
			IsAvailable: true,
		}
	} else {
		eventQueueMap[agentId] = map[chan *AgentEvent]*AgentState{
			c: {
				Version:     agentVersion,
				OsName:      osName,
				Speed:       speed,
				IP:          ip,
				IsAvailable: true,
			},
		}
	}
}

// 禁用Channel
func DisableAgentQueue(agentId string, c chan *AgentEvent) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	m, found := eventQueueMap[agentId]
	if found {
		state, ok := m[c]
		if ok {
			state.IsAvailable = false
		}
	}
}

// 删除Agent
func RemoveAgentQueue(agentId string, c chan *AgentEvent) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	_, ok := eventQueueMap[agentId]
	if ok {
		delete(eventQueueMap[agentId], c)

		if len(eventQueueMap[agentId]) == 0 {
			delete(eventQueueMap, agentId)
		}
	}
}

// 发送Agent事件
func PostAgentEvent(agentId string, event *AgentEvent) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	m, found := eventQueueMap[agentId]
	if found {
		for c, state := range m {
			if state.IsAvailable {
				c <- event
			}
		}
	}
}

// 检查Agent是否正在运行
func CheckAgentIsWaiting(agentId string) (state *AgentState, isWaiting bool) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	queue, _ := eventQueueMap[agentId]
	if len(queue) > 0 {
		for _, v := range queue {
			return v, true
		}
	}
	return nil, false
}
