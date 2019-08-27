package agentutils

// Agent状态
type AgentState struct {
	Version     string  // 版本号
	OsName      string  // 操作系统
	Speed       float64 // 连接速度，ms
	IP          string  // IP地址
	IsAvailable bool    // 是否可用
}
