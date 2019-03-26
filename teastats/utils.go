package teastats

import (
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"reflect"
	"strings"
	"sync"
)

var AllStatFilters = []maps.Map{}
var registerFilterLocker = sync.Mutex{}

var AllStartedServers = maps.Map{} // serverId => *ServerQueue
var serverLocker = sync.Mutex{}

// 注册一个或多个筛选器
func RegisterFilter(filters ...FilterInterface) {
	registerFilterLocker.Lock()
	defer registerFilterLocker.Unlock()

	for _, filter := range filters {
		for _, code := range filter.Codes() {
			periodName := ""
			dotIndex := strings.LastIndex(code, ".")
			if dotIndex > 0 {
				switch code[dotIndex+1:] {
				case ValuePeriodSecond:
					periodName = "秒"
				case ValuePeriodMinute:
					periodName = "分钟"
				case ValuePeriodHour:
					periodName = "小时"
				case ValuePeriodDay:
					periodName = "天"
				case ValuePeriodWeek:
					periodName = "周"
				case ValuePeriodMonth:
					periodName = "月"
				case ValuePeriodYear:
					periodName = "年"
				}
			}
			m := maps.Map{
				"name":     filter.Name(),
				"code":     code,
				"period":   periodName,
				"instance": filter,
			}
			AllStatFilters = append(AllStatFilters, m)
		}
	}
}

// 获取所有的filter
func FindAllStatFilters() []maps.Map {
	result := []maps.Map{}
	for _, f := range AllStatFilters {
		result = append(result, maps.Map{
			"name":   f["name"],
			"code":   f["code"],
			"period": f["period"],
		})
	}
	return result
}

// 启动一个服务的筛选器
func RestartServerFilters(serverId string, codes []string) {
	serverLocker.Lock()
	defer serverLocker.Unlock()

	// 停止现有的
	serverQueue, found := AllStartedServers[serverId]
	if found {
		serverQueue.(*ServerQueue).Stop()
	}

	// 如果没有任何指标，则删除
	if len(codes) == 0 {
		delete(AllStartedServers, serverId)
		return
	}

	queue := NewQueue()
	queue.ServerId = serverId
	queue.Start(serverId)
	serverQueue = &ServerQueue{
		Queue:   queue,
		Filters: map[string]FilterInterface{},
	}

	for _, code := range codes {
		serverQueue.(*ServerQueue).StartFilter(code)
	}

	AllStartedServers[serverId] = serverQueue
}

// 查找ServerQueue
func FindServerQueue(serverId string) *ServerQueue {
	serverQueue, found := AllStartedServers[serverId]
	if found {
		return serverQueue.(*ServerQueue)
	}
	return nil
}

// 查找单个Filter
func FindFilter(code string) FilterInterface {
	registerFilterLocker.Lock()
	defer registerFilterLocker.Unlock()
	for _, m := range AllStatFilters {
		instance := m["instance"]
		if lists.ContainsString(instance.(FilterInterface).Codes(), code) {
			return reflect.New(reflect.Indirect(reflect.ValueOf(instance)).Type()).Interface().(FilterInterface)
		}
	}
	return nil
}
