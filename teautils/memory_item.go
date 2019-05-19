package teautils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"github.com/iwind/TeaGo/logs"
	"math/big"
	"sync/atomic"
)

type MemoryItemType = int

const (
	MemoryItemTypeInt64 = 1
	MemoryItemTypeBytes = 2
)

func MemoryHashKey(key string) int64 {
	h := md5.New()
	h.Write([]byte(key))

	bi := big.NewInt(0)
	bi.SetBytes(h.Sum(nil))
	return bi.Int64()
}

type MemoryItem struct {
	Key          string
	ExpireAt     int64
	Type         MemoryItemType
	ValueInt64   int64
	ValueBytes   []byte
	IsCompressed bool
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

func (this *MemoryItem) Bytes() []byte {
	if this.IsCompressed {
		reader, err := gzip.NewReader(bytes.NewBuffer(this.ValueBytes))
		if err != nil {
			logs.Error(err)
			return this.ValueBytes
		}

		buf := make([]byte, 256)
		dataBuf := bytes.NewBuffer([]byte{})
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				dataBuf.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		return dataBuf.Bytes()
	}
	return this.ValueBytes
}

func (this *MemoryItem) String() string {
	return string(this.Bytes())
}
