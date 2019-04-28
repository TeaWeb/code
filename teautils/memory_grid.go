package teautils

import (
	"github.com/iwind/TeaGo/timers"
	"time"
)

// 内存缓存
// |           Grid               |
// | cell1, cell2, ..., cell100000 |
type MemoryGrid struct {
	cells      []*MemoryCell
	countCells int64

	recycleIndex  int
	recycleLooper *timers.Looper
}

func NewMemoryGrid(countCells int) *MemoryGrid {
	cells := []*MemoryCell{}
	if countCells <= 0 {
		countCells = 1024
	} else if countCells > 100*10000 {
		countCells = 100 * 10000
	}
	for i := 0; i < countCells; i ++ {
		cells = append(cells, NewMemoryCell())
	}

	grid := &MemoryGrid{
		cells:        cells,
		countCells:   int64(len(cells)),
		recycleIndex: -1,
	}

	grid.recycleTimer()
	return grid
}

func (this *MemoryGrid) WriteItem(item *MemoryItem) {
	if this.countCells <= 0 {
		return
	}
	hashKey := item.HashKey()
	this.cellForHashKey(hashKey).Write(hashKey, item)
}

func (this *MemoryGrid) WriteInt64(key string, value int64, lifeSeconds int64) {
	this.WriteItem(&MemoryItem{
		Key:        key,
		Type:       MemoryItemTypeInt64,
		ValueInt64: value,
		ExpireAt:   time.Now().Unix() + lifeSeconds,
	})
}

func (this *MemoryGrid) IncreaseInt64(key string, delta int64, lifeSeconds int64) (result int64) {
	hashKey := MemoryHashKey(key)
	return this.cellForHashKey(hashKey).Increase64(key, time.Now().Unix()+lifeSeconds, hashKey, delta)
}

func (this *MemoryGrid) WriteString(key string, value string, lifeSeconds int64) {
	this.WriteItem(&MemoryItem{
		Key:         key,
		Type:        MemoryItemTypeString,
		ValueString: value,
		ExpireAt:    time.Now().Unix() + lifeSeconds,
	})
}

func (this *MemoryGrid) WriteBytes(key string, value []byte, lifeSeconds int64) {
	this.WriteItem(&MemoryItem{
		Key:        key,
		Type:       MemoryItemTypeBytes,
		ValueBytes: value,
		ExpireAt:   time.Now().Unix() + lifeSeconds,
	})
}

func (this *MemoryGrid) Read(key string) *MemoryItem {
	if this.countCells <= 0 {
		return nil
	}
	hashKey := MemoryHashKey(key)
	return this.cellForHashKey(hashKey).Read(hashKey)
}

func (this *MemoryGrid) Delete(key string) {
	if this.countCells <= 0 {
		return
	}
	hashKey := MemoryHashKey(key)
	this.cellForHashKey(hashKey).Delete(hashKey)
}

func (this *MemoryGrid) Destroy() {
	if this.recycleLooper != nil {
		this.recycleLooper.Stop()
		this.recycleLooper = nil
	}
	this.cells = nil
}

func (this *MemoryGrid) cellForHashKey(hashKey int64) *MemoryCell {
	if hashKey < 0 {
		return this.cells[-hashKey%this.countCells]
	} else {
		return this.cells[hashKey%this.countCells]
	}
}

func (this *MemoryGrid) recycleTimer() {
	this.recycleLooper = timers.Loop(1*time.Minute, func(looper *timers.Looper) {
		if this.countCells == 0 {
			return
		}
		this.recycleIndex ++
		if this.recycleIndex > int(this.countCells-1) {
			this.recycleIndex = 0
		}
		this.cells[this.recycleIndex].Recycle()
	})
}
