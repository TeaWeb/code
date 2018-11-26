package teaconfigs

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
)

// 服务后端配置
type ServerBackendConfig struct {
	On            bool                   `yaml:"on" json:"on"`                       // 是否启用 @TODO
	Id            string                 `yaml:"id" json:"id"`                       // @TODO
	Name          []string               `yaml:"name" json:"name"`                   // 名称
	Address       string                 `yaml:"address" json:"address"`             // 地址
	Weight        uint                   `yaml:"weight" json:"weight"`               //@TODO
	IsBackup      bool                   `yaml:"backup" json:"isBackup"`             //@TODO
	FailTimeout   string                 `yaml:"failTimeout" json:"failTimeout"`     //@TODO
	SlowStart     string                 `yaml:"slowStart" json:"slowStart"`         //@TODO
	MaxFails      uint                   `yaml:"maxFails" json:"maxFails"`           //@TODO
	MaxConns      uint                   `yaml:"maxConns" json:"maxConns"`           //@TODO
	IsDown        bool                   `yaml:"down" json:"isDown"`                 //@TODO
	Headers       []*shared.HeaderConfig `yaml:"headers" json:"headers"`             // 自定义Header @TODO
	IgnoreHeaders []string               `yaml:"ignoreHeaders" json:"ignoreHeaders"` // 忽略的Header @TODO
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

	// Headers
	for _, header := range this.Headers {
		err := header.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// 设置Header
func (this *ServerBackendConfig) SetHeader(name string, value string) {
	found := false
	upperName := strings.ToUpper(name)
	for _, header := range this.Headers {
		if strings.ToUpper(header.Name) == upperName {
			found = true
			header.Value = value
		}
	}
	if found {
		return
	}

	header := shared.NewHeaderConfig()
	header.Name = name
	header.Value = value
	this.Headers = append(this.Headers, header)
}

// 删除指定位置上的Header
func (this *ServerBackendConfig) DeleteHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.Headers) {
		this.Headers = lists.Remove(this.Headers, index).([]*shared.HeaderConfig)
	}
}

// 取得指定位置上的Header
func (this *ServerBackendConfig) HeaderAtIndex(index int) *shared.HeaderConfig {
	if index >= 0 && index < len(this.Headers) {
		return this.Headers[index]
	}
	return nil
}

// 屏蔽一个Header
func (this *ServerBackendConfig) AddIgnoreHeader(name string) {
	this.IgnoreHeaders = append(this.IgnoreHeaders, name)
}

// 移除对Header的屏蔽
func (this *ServerBackendConfig) DeleteIgnoreHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders = lists.Remove(this.IgnoreHeaders, index).([]string)
	}
}

// 更改Header的屏蔽
func (this *ServerBackendConfig) UpdateIgnoreHeaderAtIndex(index int, name string) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders[index] = name
	}
}
