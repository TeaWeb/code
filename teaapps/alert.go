package teaapps

// 监控级别
type AlertLevel = string

const (
	AlertLevelInfo    = "info"
	AlertLevelDebug   = "debug"
	AlertLevelWarning = "warning"
	AlertLevelError   = "error"
)

// 监控设置
type Monitor struct {
	ProcessPids []uint32  // 监控的进程PID
	Files       []string  // 监控的文件
	EventTypes  []string  // 监控的事件
	Sockets     []*Socket // 监控的Socket端口
	URLs        []string  // 监控的URL
	PingHosts   []string  // 通过ICMP监控
	Scripts     []string  // 监控脚本

	Timeout  float64 // 超时时间，单位为秒
	Interval float64 // 间隔时间
	MaxFails uint32  // 最大失败次数，在达到此失败次数后才会报警

	Level AlertLevel // 级别

	Targets []*AlertTarget // 报警发送目标
}

// 报警发送目标
type AlertTarget interface {
	SetOptions(options map[string]interface{}) error // 设置选项
	Send(message string) error                       // 发送消息
}
