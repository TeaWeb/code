package teacache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"time"
)

// 内存缓存管理器
type RedisManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	Network  string
	Host     string
	Port     int
	Password string
	Sock     string

	client *redis.Client
}

func NewRedisManager() *RedisManager {
	m := &RedisManager{}
	return m
}

func (this *RedisManager) SetOptions(options map[string]interface{}) {
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}

	m := maps.NewMap(options)
	this.Network = m.GetString("network")
	this.Host = m.GetString("host")
	this.Port = m.GetInt("port")
	this.Password = m.GetString("password")
	this.Sock = m.GetString("sock")

	addr := ""
	if this.Network == "tcp" {
		if this.Port > 0 {
			addr = fmt.Sprintf("%s:%d", this.Host, this.Port)
		} else {
			addr = this.Host + ":6379"
		}
	} else if this.Network == "sock" {
		addr = this.Sock
	}

	if this.client != nil {
		this.client.Close()
	}

	this.client = redis.NewClient(&redis.Options{
		Network:      this.Network,
		Addr:         addr,
		Password:     this.Password,
		DialTimeout:  10 * time.Second, // TODO 换成可配置
		ReadTimeout:  10 * time.Second, // TODO 换成可配置
		WriteTimeout: 10 * time.Second, // TODO 换成可配置
		TLSConfig:    nil,              // TODO 支持TLS
	})
}

func (this *RedisManager) Write(key string, data []byte) error {
	cmd := this.client.Set("TEA_CACHE_"+this.id+key, string(data), this.Life)
	return cmd.Err()
}

func (this *RedisManager) Read(key string) (data []byte, err error) {
	cmd := this.client.Get("TEA_CACHE_" + this.id + key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return nil, ErrNotFound
		}
		logs.Printf("%#v", cmd.Err())
		return nil, cmd.Err()
	}
	return []byte(cmd.Val()), nil
}

// 删除
func (this *RedisManager) Delete(key string) error {
	cmd := this.client.Del("TEA_CACHE_" + this.id + key)
	return cmd.Err()
}

// 统计
func (this *RedisManager) Stat() (size int64, countKeys int, err error) {
	scan := this.client.Scan(0, "TEA_CACHE_"+this.Id()+"*", 100000)
	if scan == nil {
		return
	}
	if scan.Err() != nil {
		err = scan.Err()
		return
	}
	it := scan.Iterator()
	if it.Err() != nil {
		err = it.Err()
		return
	}
	for it.Next() {
		key := it.Val()
		b, err := this.client.Get(key).Bytes()
		if err != nil {
			continue
		}
		countKeys ++
		size += int64(len(b) + len(key))
	}
	return
}

// 清理
func (this *RedisManager) Clean() error {
	scan := this.client.Scan(0, "TEA_CACHE_"+this.Id()+"*", 100000)
	if scan == nil {
		return nil
	}
	if scan.Err() != nil {
		return scan.Err()
	}
	it := scan.Iterator()
	if it.Err() != nil {
		return it.Err()
	}
	keys := []string{}
	for it.Next() {
		key := it.Val()
		keys = append(keys, key)
	}

	if len(keys) > 0 {
		for _, key := range keys {
			this.client.Del(key)
		}
	}

	return nil
}

// 关闭
func (this *RedisManager) Close() error {
	if this.client != nil {
		//logs.Println("[cache]close cache policy instance: redis")

		err := this.client.Close()
		this.client = nil

		return err
	}

	return nil
}
