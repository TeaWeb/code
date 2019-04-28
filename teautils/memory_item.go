package teautils

import (
	"crypto/md5"
	"math/big"
	"sync/atomic"
)

type MemoryItemType = int

const (
	MemoryItemTypeInt64  = 1
	MemoryItemTypeString = 2
	MemoryItemTypeBytes  = 3
)

func MemoryHashKey(key string) int64 {
	h := md5.New()
	h.Write([]byte(key))

	bi := big.NewInt(0)
	bi.SetBytes(h.Sum(nil))
	return bi.Int64()
}

type MemoryItem struct {
	Key         string
	ExpireAt    int64
	Type        MemoryItemType
	ValueInt64  int64
	ValueString string
	ValueBytes  []byte
}

func NewMemoryItem(key string, dataType MemoryItemType) *MemoryItem {
	return &MemoryItem{
		Key:  key,
		Type: dataType,
	}
}

func (this *MemoryItem) HashKey() int64 {
	return MemoryHashKey(this.Key)
}

func (this *MemoryItem) IncreaseInt64(delta int64) {
	atomic.AddInt64(&this.ValueInt64, delta)
}
