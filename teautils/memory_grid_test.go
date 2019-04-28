package teautils

import (
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"runtime"
	"testing"
	"time"
)

func TestMemoryGrid_Write(t *testing.T) {
	grid := NewMemoryGrid(5)
	t.Log("123456:", grid.Read("123456"))

	grid.WriteInt64("abc", 1, 5)
	logs.PrintAsJSON(grid.Read("abc"), t)

	grid.WriteString("abc", "123", 5)
	logs.PrintAsJSON(grid.Read("abc"), t)

	grid.WriteBytes("abc", []byte("123"), 5)
	logs.PrintAsJSON(grid.Read("abc"), t)

	grid.Delete("abc")
	logs.PrintAsJSON(grid.Read("abc"), t)

	grid.Delete("abcd")
	grid.WriteInt64("abcd", 123, 5)

	for index, cell := range grid.cells {
		t.Log("cell:", index, len(cell.mapping), "items")
	}

	time.Sleep(10 * time.Second)
	t.Log("after recycle:")
	for index, cell := range grid.cells {
		t.Log("cell:", index, len(cell.mapping), "items")
	}

	grid.Destroy()
}

func BenchmarkMemoryGrid_Performance(b *testing.B) {
	grid := NewMemoryGrid(1024)
	for i := 0; i < b.N; i ++ {
		grid.WriteInt64(fmt.Sprintf("key:%d", i), int64(i), 3600)
	}
}

func TestMemoryGrid_Performance(t *testing.T) {
	runtime.GOMAXPROCS(1)
	now := time.Now()

	grid := NewMemoryGrid(1024)
	for i := 0; i < 1000000; i ++ {
		grid.WriteInt64(fmt.Sprintf("key:%d", i), int64(i), 3600)
		//grid.Read(fmt.Sprintf("key:%d", i-100), )
	}

	countItems := 0
	for _, cell := range grid.cells {
		countItems += len(cell.mapping)
	}
	t.Log(countItems, "items")

	t.Log(time.Since(now).Seconds()*1000, "ms")
}

func TestMemoryGrid_IncreaseInt64(t *testing.T) {
	grid := NewMemoryGrid(1024)
	grid.WriteInt64("abc", 123, 10)
	grid.IncreaseInt64("abc", 123, 10)
	grid.IncreaseInt64("abc", 123, 10)
	t.Log(grid.Read("abc"))
}

func TestMemoryGrid_Destroy(t *testing.T) {
	grid := NewMemoryGrid(1024)
	grid.WriteInt64("abc", 123, 10)
	t.Log(grid.recycleLooper, grid.cells)
	grid.Destroy()
	t.Log(grid.recycleLooper, grid.cells)
}
