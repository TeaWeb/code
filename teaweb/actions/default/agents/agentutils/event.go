package agentutils

import (
	"sync"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

var eventQueueMap = map[string]map[chan *Event]bool{} // agentId => []chan => bool
var eventQueueLocker = sync.Mutex{}

// 新Agent事件
func NewAgentEvent(name string, data interface{}) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

// 等待Agent事件
func WaitAgentQueue(agentId string, c chan *Event) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	_, ok := eventQueueMap[agentId]
	if ok {
		eventQueueMap[agentId][c] = true
	} else {
		eventQueueMap[agentId] = map[chan *Event]bool{
			c: true,
		}
	}
}

// 删除Agent
func RemoveAgentQueue(agentId string, c chan *Event) {
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
func PostAgentEvent(agentId string, event *Event) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	m, found := eventQueueMap[agentId]
	if found {
		for c, _ := range m {
			c <- event
		}
	}
}

// 是否正在运行
func CheckAgentIsWaiting(agentId string) bool {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	queue, found := eventQueueMap[agentId]
	return found && len(queue) > 0
}
