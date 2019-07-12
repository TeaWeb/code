package teaconfigs

// TCP代理配置
type TCPConfig struct {
	TCPOn bool `yaml:"tcpOn" json:"tcpOn"` // 是否支持TCP
}

// 获取新对象
func NewTCPConfig() *TCPConfig {
	return &TCPConfig{
		TCPOn: true,
	}
}
