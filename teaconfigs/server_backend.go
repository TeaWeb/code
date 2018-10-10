package teaconfigs

import (
	"strings"
	"github.com/iwind/TeaGo/utils/string"
)

// 服务后端配置
type ServerBackendConfig struct {
	On          bool     `yaml:"on" json:"on"`                   // 是否启用 @TODO
	Id          string   `yaml:"id" json:"id"`                   // @TODO
	Name        []string `yaml:"name" json:"name"`               // 名称
	Address     string   `yaml:"address" json:"address"`         // 地址
	Weight      uint     `yaml:"weight" json:"weight"`           //@TODO
	IsBackup    bool     `yaml:"backup" json:"isBackup"`         //@TODO
	FailTimeout string   `yaml:"failTimeout" json:"failTimeout"` //@TODO
	SlowStart   string   `yaml:"slowStart" json:"slowStart"`     //@TODO
	MaxFails    uint     `yaml:"maxFails" json:"maxFails"`       //@TODO
	MaxConns    uint     `yaml:"maxConns" json:"maxConns"`       //@TODO
	IsDown      bool     `yaml:"down" json:"isDown"`             //@TODO
}

func NewServerBackendConfig() *ServerBackendConfig {
	return &ServerBackendConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

func (this *ServerBackendConfig) Validate() error {
	// 是否有端口
	if strings.Index(this.Address, ":") == -1 {
		// @TODO 如果是tls，则为443
		this.Address += ":80"
	}

	return nil
}
