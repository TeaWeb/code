package configs

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"sync"
	"github.com/iwind/TeaGo/logs"
	"net/url"
	"strings"
	"github.com/iwind/TeaGo/types"
	"fmt"
)

// MongoDB配置
type MongoConfig struct {
	Scheme     string `json:"scheme"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       uint   `json:"port"`
	RequestURI string // @TODO 未来版本需要实现
}

var mongoConfig *MongoConfig
var mongoConfigLocker sync.Mutex

// 取得全局的MongoDB配置
func SharedMongoConfig() *MongoConfig {
	mongoConfigLocker.Lock()
	defer mongoConfigLocker.Unlock()

	if mongoConfig != nil {
		return mongoConfig
	}

	mongoConfig = &MongoConfig{}

	confFile := Tea.ConfigFile("mongo.conf")
	reader, err := files.NewReader(confFile)
	if err != nil {
		return mongoConfig
	}
	defer reader.Close()

	m, err := reader.ReadYAMLMap()
	if err != nil {
		logs.Error(err)
		return mongoConfig
	}

	mongoURL := m.GetString("uri")
	if len(mongoURL) == 0 {
		return mongoConfig
	}

	u, err := url.Parse(mongoURL)
	if err != nil {
		logs.Error(err)
		return mongoConfig
	}

	mongoConfig.Scheme = u.Scheme

	index := strings.Index(u.Host, ":")
	if index >= 0 {
		mongoConfig.Host = u.Host[:index]
		mongoConfig.Port = types.Uint(u.Port())
	} else {
		mongoConfig.Host = u.Host
		mongoConfig.Port = types.Uint(u.Port())
	}

	if u.User != nil {
		mongoConfig.Username = u.User.Username()
		mongoConfig.Password, _ = u.User.Password()
	}

	return mongoConfig
}

// 组合后的URI
func (this *MongoConfig) URI() string {
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

// 保存修改后的MongoDB配置
func (this *MongoConfig) WriteBack() error {
	confFile := Tea.ConfigFile("mongo.conf")
	writer, err := files.NewWriter(confFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = writer.WriteYAML(map[string]interface{}{
		"uri": this.URI(),
	})
	return err
}
