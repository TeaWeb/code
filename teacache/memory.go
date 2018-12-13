package teacache

import (
	"errors"
	"github.com/iwind/TeaGo/caches"
	"sync"
	"time"
)

// 内存缓存管理器
type MemoryManager struct {
	Capacity float64       // 容量
	Life     time.Duration // 有效期

	cache        *caches.Factory
	memory       float64
	memoryLocker sync.Mutex
}

func NewMemoryManager() *MemoryManager {
	m := &MemoryManager{}

	factory := caches.NewFactory()
	factory.OnOperation(func(op caches.CacheOperation, value interface{}) {
		m.memoryLocker.Lock()
		defer m.memoryLocker.Unlock()
		if op == caches.CacheOperationSet {
			m.memory += float64(len(value.([]byte)))
		} else if op == caches.CacheOperationDelete {
			m.memory -= float64(len(value.([]byte)))
		}
	})
	m.cache = factory
	return m
}

func (this *MemoryManager) SetOptions(options map[string]interface{}) {

}

func (this *MemoryManager) Write(key string, data []byte) error {
	// 检查容量
	if this.Capacity > 0 && this.memory+float64(len(data)) >= this.Capacity {
		return errors.New("out of memory")
	}

	this.cache.Set(key, data).Expire(this.Life)
	return nil
}

func (this *MemoryManager) Read(key string) (data []byte, err error) {
	value, found := this.cache.Get(key)
	if !found {
		return nil, ErrNotFound
	}
	return value.([]byte), nil
}
