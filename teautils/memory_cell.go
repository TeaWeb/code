package teautils

import (
	"sync"
	"time"
)

type MemoryCell struct {
	mapping map[int64]*MemoryItem // key => item
	locker  sync.RWMutex
}

func NewMemoryCell() *MemoryCell {
	return &MemoryCell{
		mapping: map[int64]*MemoryItem{},
	}
}

func (this *MemoryCell) Write(hashKey int64, item *MemoryItem) {
	if item == nil {
		return
	}
	this.locker.Lock()
	this.mapping[hashKey] = item
	this.locker.Unlock()
}

func (this *MemoryCell) Increase64(key string, expireAt int64, hashKey int64, delta int64) (result int64) {
	this.locker.Lock()
	item, ok := this.mapping[hashKey]
	if ok {
		// reset to zero if expired
		if item.ExpireAt < time.Now().Unix() {
			item.ValueInt64 = 0
			item.ExpireAt = expireAt
		}
		item.IncreaseInt64(delta)
		result = item.ValueInt64
	} else {
		item := NewMemoryItem(key, MemoryItemTypeInt64)
		item.ValueInt64 = delta
		item.ExpireAt = expireAt
		this.mapping[hashKey] = item
		result = delta
	}
	this.locker.Unlock()
	return
}

func (this *MemoryCell) Read(hashKey int64) *MemoryItem {
	this.locker.RLock()

	item, ok := this.mapping[hashKey]
	if ok {
		this.locker.RUnlock()

		if item.ExpireAt < time.Now().Unix() {
			return nil
		}
		return item
	}

	this.locker.RUnlock()
	return nil
}

func (this *MemoryCell) Delete(hashKey int64) {
	this.locker.Lock()
	delete(this.mapping, hashKey)
	this.locker.Unlock()
}

func (this *MemoryCell) Recycle() {
	this.locker.Lock()
	timestamp := time.Now().Unix()
	for key, item := range this.mapping {
		if item.ExpireAt < timestamp {
			delete(this.mapping, key)
		}
	}
	this.locker.Unlock()
}
