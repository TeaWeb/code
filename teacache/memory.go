package teacache

import (
	"github.com/iwind/TeaGo/caches"
	"sync"
	"time"
)

// 内存缓存管理器
type MemoryManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	cache        *caches.Factory
	memory       int64
	count        int
	memoryLocker sync.Mutex
}

func NewMemoryManager() *MemoryManager {
	m := &MemoryManager{}

	factory := caches.NewFactory()
	factory.OnOperation(func(op caches.CacheOperation, value interface{}) {
		m.memoryLocker.Lock()
		defer m.memoryLocker.Unlock()
		if op == caches.CacheOperationSet {
			m.memory += int64(len(value.([]byte)))
			m.count ++
		} else if op == caches.CacheOperationDelete {
			m.memory -= int64(len(value.([]byte)))
			m.count --
		}
	})
	m.cache = factory

	return m
}

func (this *MemoryManager) SetOptions(options map[string]interface{}) {
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}
}

func (this *MemoryManager) Write(key string, data []byte) error {
	// 检查容量
	if this.Capacity > 0 && float64(this.memory+int64(len(data))) >= this.Capacity {
		this.memory = 0
		this.cache.Reset()
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

// 统计
func (this *MemoryManager) Stat() (size int64, countKeys int, err error) {
	return this.memory, this.count, nil
}

func (this *MemoryManager) Close() error {
	if this.cache == nil {
		return nil
	}
	//logs.Println("[cache]close cache policy instance: memory")
	this.cache.Close()
	return nil
}
