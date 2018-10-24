package teaapps

import (
	"github.com/TeaWeb/plugin/apps"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

// 接口
type AppInterface interface {
	Start() error
	Stop() error
}

// App定义
type App struct {
	Id        string `json:"id"` // 唯一ID，通常系统会自动生成
	Name      string `json:"name"`
	Developer string `json:"developer"`
	Site      string `json:"site"`
	DocSite   string `json:"docSite"`
	Version   string `json:"version"`
	Icon      []byte `json:"icon"`

	Processes  []*Process    `json:"processes"`
	Operations []*Operation  `json:"operations"`
	Monitors   []*Monitor    `json:"monitors"`
	Statistics []*Statistics `json:"statistics"`
	Logs       []*Log        `json:"logs"`

	IsRunning bool `json:"isRunning"`

	onReloadFunc func()
}

// 取得新App
func NewApp() *App {
	return &App{}
}

// 重置进程
func (this *App) ResetProcesses() {
	this.Processes = []*Process{}
}

// 添加进程，不包括子进程
func (this *App) AddProcess(process ... *Process) {
	this.Processes = append(this.Processes, process ...)
}

// 计算总体CPU占用量
func (this *App) SumCPUUsage() *CPUUsage {
	cpuUsage := &CPUUsage{}
	for _, process := range this.Processes {
		if process.CPUUsage != nil {
			cpuUsage.Percent += process.CPUUsage.Percent
		}
	}
	return cpuUsage
}

// 计算总体内存占用量
func (this *App) SumMemoryUsage() *MemoryUsage {
	memoryUsage := &MemoryUsage{}
	for _, process := range this.Processes {
		if process.MemoryUsage != nil {
			memoryUsage.Percent += process.MemoryUsage.Percent
			memoryUsage.RSS += process.MemoryUsage.RSS
			memoryUsage.VMS += process.MemoryUsage.VMS
		}
	}
	return memoryUsage
}

// 计算总连接数
func (this *App) CountAllConnections() int {
	m := maps.Map{}
	for _, process := range this.Processes {
		for _, conn := range process.Connections {
			m[conn] = true
		}
	}
	return m.Len()
}

// 计算总打开的文件数量
func (this *App) CountAllOpenFiles() int {
	m := maps.Map{}
	for _, process := range this.Processes {
		for _, file := range process.OpenFiles {
			m[file] = true
		}
	}
	return m.Len()
}

func (this *App) CountAllListens() int {
	m := maps.Map{}
	for _, process := range this.Processes {
		for _, listen := range process.Listens {
			m[listen.Network+"://"+listen.Addr] = true
		}
	}
	return m.Len()
}

func (this *App) LoadFromInterface(a *apps.App) {
	this.Id = a.Id
	this.Name = a.Name
	this.Developer = a.Developer
	this.Site = a.Site
	this.DocSite = a.DocSite
	this.Version = a.Version
	this.Icon = a.Icon
	this.IsRunning = a.IsRunning

	// Processes
	this.Processes = []*Process{}

	for _, p := range a.Processes {
		p2 := NewProcess()
		p2.IsRunning = p.IsRunning
		p2.Pid = p.Pid
		p2.Name = p.Name
		p2.Ppid = p.Ppid
		p2.Cwd = p.Cwd
		p2.User = p.User
		p2.Uid = p.Uid
		p2.Gid = p.Gid
		p2.CreateTime = p.CreateTime
		p2.Cmdline = p.Cmdline
		p2.File = p.File
		p2.Dir = p.Dir
		if p.CPUUsage != nil {
			p2.CPUUsage = &CPUUsage{
				Percent: p.CPUUsage.Percent,
			}
		}
		if p.MemoryUsage != nil {
			p2.MemoryUsage = &MemoryUsage{
				RSS:     p.MemoryUsage.RSS,
				VMS:     p.MemoryUsage.VMS,
				Percent: p.MemoryUsage.Percent,
			}
		}
		if p.OpenFiles != nil {
			p2.OpenFiles = p.OpenFiles
		} else {
			p2.OpenFiles = []string{}
		}
		if p.Connections != nil {
			p2.Connections = p.Connections
		} else {
			p2.Connections = []string{}
		}

		if len(p.Listens) > 0 {
			for _, listen := range p.Listens {
				p2.Listens = append(p2.Listens, &Listen{
					Network: listen.Network,
					Addr:    listen.Addr,
				})
			}
		} else {
			p2.Listens = []*Listen{}
		}
		this.AddProcess(p2)
	}

	// Operations TODO
	this.Operations = []*Operation{}

	// Monitors TODO
	this.Monitors = []*Monitor{}

	// Statistics TODO
	this.Statistics = []*Statistics{}

	// Logs TODO
	this.Logs = []*Log{}
}

func (this *App) OnReload(f func()) {
	this.onReloadFunc = f
}

func (this *App) Reload() {
	if this.onReloadFunc != nil {
		this.onReloadFunc()
	}
}

// 获取UniqueId
func (this *App) UniqueId() string {
	if len(this.Processes) == 0 {
		return stringutil.Md5(this.Name + "@" + this.Version)
	}
	return stringutil.Md5(this.Processes[0].File + "@" + this.Processes[0].Name)
}
