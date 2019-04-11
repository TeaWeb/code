package teacache

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/utils/string"
	"time"
)

var ErrNotFound = errors.New("cache not found")

// 缓存管理接口
type ManagerInterface interface {
	// 写入
	Write(key string, data []byte) error

	// 读取
	Read(key string) (data []byte, err error)

	// 设置选项
	SetOptions(options map[string]interface{})

	// 关闭
	Close() error
}

// 获取新的管理对象
func NewManagerFromConfig(config *shared.CachePolicy) ManagerInterface {
	switch config.Type {
	case "memory":
		m := NewMemoryManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		return m
	case "file":
		m := NewFileManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		return m
	case "redis":
		m := NewRedisManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		return m
	case "leveldb":
		m := NewLevelDBManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		return m
	}
	return nil
}
