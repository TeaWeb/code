package teautils

import (
	"sync"
	"testing"
)

func TestMemoryItem_HashKey(t *testing.T) {
	{
		item := NewMemoryItem("", MemoryItemTypeBytes)
		t.Log(item.HashKey())
	}

	{
		item := NewMemoryItem("0", MemoryItemTypeBytes)
		t.Log(item.HashKey())
	}

	{
		item := NewMemoryItem("123", MemoryItemTypeBytes)
		t.Log(item.HashKey())
	}

	{
		item := NewMemoryItem("456", MemoryItemTypeBytes)
		t.Log(item.HashKey())
	}

	{
		item := NewMemoryItem("123456", MemoryItemTypeBytes)
		t.Log(item.HashKey())
	}
}

func TestMemoryItem_Increase(t *testing.T) {
	item := NewMemoryItem("123456", MemoryItemTypeInt64)
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i ++ {
		go func() {
			item.IncreaseInt64(1)
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log(item.ValueInt64)
	item.IncreaseInt64(-100)
	t.Log(item.ValueInt64)
}
