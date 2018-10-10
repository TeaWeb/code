package teaproxy

import (
	"runtime"
	"github.com/pbnjay/memory"
	"sync"
	"time"
)

type FixedCache struct {
	maxMemory uint64
	items     map[string]interface{}
	mutex     *sync.Mutex
}

func NewFixedCache() *FixedCache {
	cache := &FixedCache{
		maxMemory: memory.TotalMemory() / 8,
		items:     map[string]interface{}{},
		mutex:     &sync.Mutex{},
	}

	go cache.checkMemory()

	return cache
}

func (this *FixedCache) Add(key string, object interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.items[key] = object
}

func (this *FixedCache) Get(key string) (interface{}, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	object, found := this.items[key]
	return object, found
}

func (this *FixedCache) Trim() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	max := len(this.items) / 2
	count := 0
	for key := range this.items {
		if count < max {
			delete(this.items, key)
		} else {
			break
		}

		count ++
	}
}

func (this *FixedCache) checkMemory() {
	for {
		func() {
			stat := &runtime.MemStats{}
			runtime.ReadMemStats(stat)

			total := stat.TotalAlloc
			if total > this.maxMemory {
				this.Trim()
			}
		}()
		time.Sleep(5 * time.Second)
	}
}
