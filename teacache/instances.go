package teacache

import (
	"sync"
)

var cachePolicyMap = map[string]ManagerInterface{}
var cachePolicyMapLocker = sync.RWMutex{}

func ResetCachePolicyManager(filename string) {
	cachePolicyMapLocker.Lock()
	defer cachePolicyMapLocker.Unlock()

	manager, ok := cachePolicyMap[filename]
	if ok {
		manager.Close()
		delete(cachePolicyMap, filename)
	}
}
