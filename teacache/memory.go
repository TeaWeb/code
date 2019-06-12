package teacache

import (
	"errors"
	"github.com/TeaWeb/code/teamemory"
	"time"
)

// 内存缓存管理器
type MemoryManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	grid *teamemory.Grid
}

func NewMemoryManager() *MemoryManager {
	m := &MemoryManager{}

	return m
}

func (this *MemoryManager) SetOptions(options map[string]interface{}) {
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}

	countCells := 128
	opts := []interface{}{}
	if this.Capacity > 0 {
		capacityBytes := int64(this.Capacity)
		opts = append(opts, teamemory.NewLimitSizeOpt(capacityBytes))
		countCells = int(capacityBytes / 1024 / 1024 / 128)
	}
	this.grid = teamemory.NewGrid(countCells, opts ...)
}

func (this *MemoryManager) Write(key string, data []byte) error {
	if this.grid == nil {
		return errors.New("grid has not been initialized")
	}
	life := int64(this.Life.Seconds())
	this.grid.WriteBytes([]byte(key), data, life)
	return nil
}

func (this *MemoryManager) Read(key string) (data []byte, err error) {
	if this.grid == nil {
		return nil, errors.New("grid has not been initialized")
	}
	item := this.grid.Read([]byte(key))
	if item == nil {
		return nil, ErrNotFound
	}
	return item.Bytes(), nil
}

// 删除
func (this *MemoryManager) Delete(key string) error {
	this.grid.Delete([]byte(key))
	return nil
}

// 统计
func (this *MemoryManager) Stat() (size int64, countKeys int, err error) {
	stat := this.grid.Stat()
	return stat.TotalBytes, stat.CountItems, nil
}

// 清理
func (this *MemoryManager) Clean() error {
	this.grid.Reset()
	return nil
}

// 关闭
func (this *MemoryManager) Close() error {
	if this.grid == nil {
		return nil
	}
	//logs.Println("[cache]close cache policy instance: memory")
	this.grid.Destroy()
	return nil
}
