package configs

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"net"
	"net/url"
	"strings"
	"sync"
)

// MongoDB连接配置
type MongoConnectionConfig struct {
	Scheme                  string             `json:"scheme"`
	Username                string             `json:"username"`
	Password                string             `json:"password"`
	Host                    string             `json:"host"`
	Port                    uint               `json:"port"`
	AuthMechanism           string             `json:"authMechanism"`
	AuthMechanismProperties []*shared.Variable `json:"authMechanismProperties"`
	RequestURI              string             `json:"requestURI"` // @TODO 未来版本需要实现
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

	host, port, err := net.SplitHostPort(u.Host)
	if err == nil {
		mongoConnectionConfig.Host = host
		mongoConnectionConfig.Port = types.Uint(port)
	} else {
		mongoConnectionConfig.Host = u.Host
		mongoConnectionConfig.Port = types.Uint(u.Port())
	}

	if u.User != nil {
		mongoConnectionConfig.Username = u.User.Username()
		mongoConnectionConfig.Password, _ = u.User.Password()
	}

	mongoConnectionConfig.AuthMechanism = u.Query().Get("authMechanism")
	properties := u.Query().Get("authMechanismProperties")
	if len(properties) > 0 {
		for _, property := range strings.Split(properties, ",") {
			if strings.Contains(property, ":") {
				pieces := strings.Split(property, ":")
				mongoConnectionConfig.AuthMechanismProperties = append(mongoConnectionConfig.AuthMechanismProperties, shared.NewVariable(pieces[0], pieces[1]))
			}
		}
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

	if len(this.AuthMechanism) > 0 {
		uri += "/?authMechanism=" + this.AuthMechanism

		if len(this.AuthMechanismProperties) > 0 {
			properties := []string{}
			for _, v := range this.AuthMechanismProperties {
				properties = append(properties, v.Name+":"+v.Value)
			}
			uri += "&authMechanismProperties=" + strings.Join(properties, ",")
		}
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

	if len(this.AuthMechanism) > 0 {
		uri += "/?authMechanism=" + this.AuthMechanism

		if len(this.AuthMechanismProperties) > 0 {
			properties := []string{}
			for _, v := range this.AuthMechanismProperties {
				properties = append(properties, v.Name+":"+v.Value)
			}
			uri += "&authMechanismProperties=" + strings.Join(properties, ",")
		}
	}

	return uri
}

// 取得Map形式的认证属性
func (this *MongoConnectionConfig) AuthMechanismPropertiesMap() map[string]string {
	m := map[string]string{}
	for _, v := range this.AuthMechanismProperties {
		m[v.Name] = v.Value
	}
	return m
}

// 取得字符串形式的认证属性
func (this *MongoConnectionConfig) AuthMechanismPropertiesString() string {
	s := []string{}
	for _, v := range this.AuthMechanismProperties {
		s = append(s, v.Name+":"+v.Value)
	}
	return strings.Join(s, ",")
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
