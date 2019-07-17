package teaconfigs

// TCP代理配置
type TCPConfig struct {
	TCPOn         bool `yaml:"tcpOn" json:"tcpOn"`                 // 是否开启TCP
	FailReconnect bool `yaml:"failReconnect" json:"failReconnect"` // 失败是否重连
	FailResend    bool `yaml:"failResend" json:"failResend"`       // 失败是否重发
}

// 获取新对象
func NewTCPConfig() *TCPConfig {
	return &TCPConfig{
		TCPOn: true,
	}
}

// 校验
func (this *TCPConfig) Validate() error {
	return nil
}
