package agentutils

import (
	"sync"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Agent状态
type State struct {
	Version string  // 版本号
	OsName  string  // 操作系统
	Speed   float64 // 连接速度，ms
	IP      string  // IP地址
}

var eventQueueMap = map[string]map[chan *Event]*State{} // agentId => { chan => State }
var eventQueueLocker = sync.Mutex{}

// 新Agent事件
func NewAgentEvent(name string, data interface{}) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

// 等待Agent事件
func WaitAgentQueue(agentId string, agentVersion string, osName string, speed float64, ip string, c chan *Event) {
	eventQueueLocker.Lock()
	defer eventQueueLocker.Unlock()
	_, found := eventQueueMap[agentId]
	if found {
		eventQueueMap[agentId][c] = &State{
			Version: agentVersion,
			OsName:  osName,
			Speed:   speed,
			IP:      ip,
		}
	} else {
		eventQueueMap[agentId] = map[chan *Event]*State{
			c: {
				Version: agentVersion,
				OsName:  osName,
				Speed:   speed,
				IP:      ip,
			},
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
func CheckAgentIsWaiting(agentId string) (state *State, isWaiting bool) {
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
