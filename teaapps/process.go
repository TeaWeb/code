package teaapps

type Listen struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
}

type Process struct {
	Name string `json:"name"`
	Pid  int32  `json:"pid"`
	Ppid int32  `json:"ppid"`
	Cwd  string `json:"cwd"`

	User string `json:"user"`
	Uid  int32  `json:"uid"`
	Gid  int32  `json:"gid"`

	CreateTime  int64        `json:"createTime"` // 时间戳
	Cmdline     string       `json:"cmdline"`    //命令行
	File        string       `json:"file"`       // 命令行文件路径
	Dir         string       `json:"dir"`        // 命令行文件所在目录
	CPUUsage    *CPUUsage    `json:"cpuUsage"`
	MemoryUsage *MemoryUsage `json:"memoryUsage"`

	OpenFiles   []string  `json:"openFiles"`
	Connections []string  `json:"connections"`
	Listens     []*Listen `json:"listens"`

	IsRunning bool `json:"isRunning"`
}

func NewProcess() *Process {
	return &Process{}
}
