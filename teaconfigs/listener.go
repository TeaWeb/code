package teaconfigs

import (
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
	"errors"
	"github.com/go-yaml/yaml"
	"path/filepath"
	"strings"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
)

// 本地监听服务配置
type ListenerConfig struct {
	Address string
	Http    bool
	SSL     *SSLConfig
	Servers []*ServerConfig
}

// 从配置文件中分析配置
func ParseConfigs() ([]*ListenerConfig, error) {
	listenerConfigMap := map[string]*ListenerConfig{}

	configsDir := Tea.ConfigDir()
	files, err := filepath.Glob(configsDir + Tea.DS + "*.proxy.conf")
	if err != nil {
		return nil, err
	}

	for _, configFile := range files {
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			logs.Error(err)
			continue
		}

		serverConfig := &ServerConfig{}
		err = yaml.Unmarshal(configData, serverConfig)
		if err != nil {
			logs.Error(err)
			continue
		}

		if len(serverConfig.Listen) == 0 {
			logs.Error(errors.New("'listen' in config should not be empty"))
			continue
		}

		err = serverConfig.Validate()
		if err != nil {
			logs.Error(err)
			continue
		}

		if !serverConfig.On {
			continue
		}

		for _, address := range serverConfig.Listen {
			// 是否有端口
			if strings.Index(address, ":") == -1 {
				if serverConfig.SSL != nil && serverConfig.SSL.On {
					address += ":443"
				} else {
					address += ":80"
				}
			}

			listenerConfig, found := listenerConfigMap[address]
			if !found {
				listenerConfig = &ListenerConfig{
					Address: address,
					Servers: []*ServerConfig{serverConfig},
				}
				listenerConfigMap[address] = listenerConfig
			} else {
				listenerConfig.Servers = append(listenerConfig.Servers, serverConfig)
			}

			if serverConfig.SSL != nil {
				listenerConfig.SSL = serverConfig.SSL
			}
			if serverConfig.Http {
				listenerConfig.Http = true
			}
		}
	}

	listenerConfigArray := []*ListenerConfig{}
	for _, listenerConfig := range listenerConfigMap {
		listenerConfigArray = append(listenerConfigArray, listenerConfig)
	}

	return listenerConfigArray, nil
}

// 获取当前监听服务的端口
func (this *ListenerConfig) Port() int {
	index := strings.LastIndex(this.Address, ":")
	if index < 0 {
		return 0
	}
	return types.Int(this.Address[index+1:])
}

// 添加服务
func (this *ListenerConfig) AddServer(serverConfig *ServerConfig) {
	this.Servers = append(this.Servers, serverConfig)
}

// 根据域名来查找匹配的域名
// @TODO 把查找的结果加入缓存
func (this *ListenerConfig) FindNamedServer(name string) (serverConfig *ServerConfig, serverName string) {
	countServers := len(this.Servers)
	if countServers == 0 {
		return nil, ""
	}

	// 如果只有一个server，则默认为这个
	if countServers == 1 {
		server := this.Servers[0]
		matchedName, matched := server.MatchName(name)
		if matched {
			if len(matchedName) > 0 {
				return server, matchedName
			} else {
				return server, name
			}
		}

		// 匹配第一个域名
		firstName := server.FirstName()
		if len(firstName) > 0 {
			return server, firstName
		}
		return server, name
	}

	// 精确查找
	for _, server := range this.Servers {
		if lists.Contains(server.Name, name) {
			return server, name
		}
	}

	// 模糊查找
	for _, server := range this.Servers {
		if _, matched := server.MatchName(name); matched {
			return server, name
		}
	}

	// 如果没有找到，则匹配到第一个
	server := this.Servers[0]
	firstName := server.FirstName()
	if len(firstName) > 0 {
		return server, firstName
	}

	return server, name
}
