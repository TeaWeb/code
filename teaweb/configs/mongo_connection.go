package configs

import (
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"net/url"
	"strings"
	"sync"
)

// MongoDB连接配置
type MongoConnectionConfig struct {
	Scheme     string `json:"scheme"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       uint   `json:"port"`
	RequestURI string // @TODO 未来版本需要实现
}

var mongoConnectionConfig *MongoConnectionConfig
var mongoConnectionConfigLocker sync.Mutex

// 取得全局的MongoDB配置
func SharedMongoConfig() *MongoConnectionConfig {
	mongoConnectionConfigLocker.Lock()
	defer mongoConnectionConfigLocker.Unlock()

	if mongoConnectionConfig != nil {
		return mongoConnectionConfig
	}

	mongoConnectionConfig = &MongoConnectionConfig{}

	mongoConfig, err := LoadMongoConfig()
	if err != nil {
		logs.Error(err)
		return mongoConnectionConfig
	}

	u, err := url.Parse(mongoConfig.URI)
	if err != nil {
		logs.Error(err)
		return mongoConnectionConfig
	}

	mongoConnectionConfig.Scheme = u.Scheme

	index := strings.Index(u.Host, ":")
	if index >= 0 {
		mongoConnectionConfig.Host = u.Host[:index]
		mongoConnectionConfig.Port = types.Uint(u.Port())
	} else {
		mongoConnectionConfig.Host = u.Host
		mongoConnectionConfig.Port = types.Uint(u.Port())
	}

	if u.User != nil {
		mongoConnectionConfig.Username = u.User.Username()
		mongoConnectionConfig.Password, _ = u.User.Password()
	}

	return mongoConnectionConfig
}

// 组合后的URI
func (this *MongoConnectionConfig) URI() string {
	uri := ""
	if len(this.Scheme) > 0 {
		uri += this.Scheme + "://"
	} else {
		uri += "mongodb://"
	}

	if len(this.Username) > 0 {
		uri += this.Username
		if len(this.Password) > 0 {
			uri += ":" + this.Password
		}
		uri += "@"
	}

	uri += this.Host
	if this.Port > 0 {
		uri += ":" + fmt.Sprintf("%d", this.Port)
	}

	return uri
}

// 组合后的URI，但是对URI进行掩码
func (this *MongoConnectionConfig) URIMask() string {
	uri := ""
	if len(this.Scheme) > 0 {
		uri += this.Scheme + "://"
	} else {
		uri += "mongodb://"
	}

	if len(this.Username) > 0 {
		uri += this.Username
		if len(this.Password) > 0 {
			uri += ":" + strings.Repeat("*", len(this.Password))
		}
		uri += "@"
	}

	uri += this.Host
	if this.Port > 0 {
		uri += ":" + fmt.Sprintf("%d", this.Port)
	}

	return uri
}

// 保存修改后的MongoDB配置
func (this *MongoConnectionConfig) Save() error {
	config, _ := LoadMongoConfig()
	if config == nil {
		config = NewMongoConfig()
		config.URI = this.URI()
	} else {
		config.URI = this.URI()
	}
	return config.Save()
}
