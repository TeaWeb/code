package teautils

import (
	"compress/gzip"
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

func TestMemoryGrid_Compress(t *testing.T) {
	grid := NewMemoryGrid(5, NewMemoryGridCompressOpt(1))
	grid.WriteString("hello", strings.Repeat("abcd", 10240), 30)
	t.Log(len(string(grid.Read("hello").String())))
	t.Log(len(grid.Read("hello").ValueBytes))
}

func BenchmarkMemoryGrid_Performance(b *testing.B) {
	grid := NewMemoryGrid(1024)
	for i := 0; i < b.N; i ++ {
		grid.WriteInt64("key:"+strconv.Itoa(i), int64(i), 3600)
	}
}

func TestMemoryGrid_Performance(t *testing.T) {
	runtime.GOMAXPROCS(1)

	grid := NewMemoryGrid(1024)

	now := time.Now()

	s := []byte(strings.Repeat("abcd", 10*1024))

	for i := 0; i < 100000; i ++ {
		grid.WriteBytes(fmt.Sprintf("key:%d_%d", i, 1), s, 3600)
		item := grid.Read(fmt.Sprintf("key:%d_%d", i, 1))
		if item != nil {
			_ = item.String()
		}
	}

	countItems := 0
	for _, cell := range grid.cells {
		countItems += len(cell.mapping)
	}
	t.Log(countItems, "items")

	t.Log(time.Since(now).Seconds()*1000, "ms")
}

func TestMemoryGrid_Performance_Concurrent(t *testing.T) {
	//runtime.GOMAXPROCS(1)

	grid := NewMemoryGrid(1024)

	now := time.Now()

	s := []byte(strings.Repeat("abcd", 10*1024))

	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	for c := 0; c < runtime.NumCPU(); c ++ {
		go func(c int) {
			defer wg.Done()
			for i := 0; i < 50000; i ++ {
				grid.WriteBytes(fmt.Sprintf("key:%d_%d", i, c), s, 3600)
				item := grid.Read(fmt.Sprintf("key:%d_%d", i, c))
				if item != nil {
					_ = item.String()
				}
			}
		}(c)
	}

	wg.Wait()
	countItems := 0
	for _, cell := range grid.cells {
		countItems += len(cell.mapping)
	}
	t.Log(countItems, "items")

	t.Log(time.Since(now).Seconds()*1000, "ms")
}

func TestMemoryGrid_CompressPerformance(t *testing.T) {
	runtime.GOMAXPROCS(1)

	grid := NewMemoryGrid(1024, NewMemoryGridCompressOpt(gzip.BestCompression))

	now := time.Now()
	data := []byte(strings.Repeat("abcd", 1024))

	for i := 0; i < 100000; i ++ {
		grid.WriteBytes(fmt.Sprintf("key:%d", i), data, 3600)
		item := grid.Read(fmt.Sprintf("key:%d", i+100))
		if item != nil {
			_ = item.String()
		}
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
