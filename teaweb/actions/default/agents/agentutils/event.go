package agentutils

import (
	"sync"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

var eventQueueMap = map[string]map[chan *Event]string{} // agentId => []chan => string
var eventQueueLocker = sync.Mutex{}

// 新Agent事件
func NewAgentEvent(name string, data interface{}) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

// 等待Agent事件
func WaitAgentQueue(agentId string, agentVersion string, c chan *Event) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	_, ok := eventQueueMap[agentId]
	if ok {
		eventQueueMap[agentId][c] = agentVersion
	} else {
		eventQueueMap[agentId] = map[chan *Event]string{
			c: agentVersion,
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

// 检查Agent是否正在运行
func CheckAgentIsWaiting(agentId string) (version string, isWaiting bool) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	queue, _ := eventQueueMap[agentId]
	if len(queue) > 0 {
		for _, v := range queue {
			return v, true
		}
	}
	return "", false
}
