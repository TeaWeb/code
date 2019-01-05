package teaconfigs

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/maps"
	"time"
)

// websocket设置
type WebsocketConfig struct {
	// 后端服务器列表
	BackendList `yaml:",inline"`

	On bool `yaml:"on" json:"on"` // 是否开启

	// 握手超时时间
	HandshakeTimeout string `yaml:"handshakeTimeout" json:"handshakeTimeout"`

	// 允许的域名，支持 www.example.com, example.com, .example.com, *.example.com
	AllowAllOrigins bool     `yaml:"allowAllOrigins" json:"allowAllOrigins"`
	Origins         []string `yaml:"origins" json:"origins"`

	// 转发方式
	ForwardMode WebsocketForwardMode `yaml:"forwardMode" json:"forwardMode"`

	handshakeTimeoutDuration time.Duration
}

// 获取新对象
func NewWebsocketConfig() *WebsocketConfig {
	return &WebsocketConfig{
		On: true,
	}
}

// 校验
func (this *WebsocketConfig) Validate() error {
	// backends
	err := this.ValidateBackends()
	if err != nil {
		return err
	}

	// duration
	this.handshakeTimeoutDuration, _ = time.ParseDuration(this.HandshakeTimeout)

	return nil
}

// 获取握手超时时间
func (this *WebsocketConfig) HandshakeTimeoutDuration() time.Duration {
	return this.handshakeTimeoutDuration
}

// 转发模式名称
func (this *WebsocketConfig) ForwardModeSummary() maps.Map {
	for _, mode := range AllWebsocketForwardModes() {
		if mode["mode"] == this.ForwardMode {
			return mode
		}
	}
	return nil
}

// 匹配域名
func (this *WebsocketConfig) MatchOrigin(origin string) bool {
	if this.AllowAllOrigins {
		return true
	}
	return teautils.MatchDomains(this.Origins, origin)
}
