package cmd

import (
	"fmt"
	"github.com/TeaWeb/code/teacluster"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var sharedShell *WebShell = nil

// 命令行相关封装
type WebShell struct {
	ShouldStop bool
}

// 获取新对象
func NewWebShell() *WebShell {
	sharedShell = &WebShell{}
	return sharedShell
}

// 启动
func (this *WebShell) Start(server *TeaGo.Server) {
	// 重置ROOT
	this.resetRoot()

	// 执行参数
	if this.execArgs(os.Stdout) {
		this.ShouldStop = true
		return
	}

	// 当前PID
	err := files.NewFile(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid").
		WriteString(fmt.Sprintf("%d", os.Getpid()))
	if err != nil {
		logs.Error(err)
	}

	// 信号
	signalsChannel := make(chan os.Signal, 16)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGHUP, syscall.Signal(0x1e /**syscall.SIGUSR1**/), syscall.Signal(0x1f /**syscall.SIGUSR2**/), syscall.SIGTERM)
	go func() {
		for {
			sig := <-signalsChannel

			if sig == syscall.SIGHUP { // 重置
				configs.SharedAdminConfig().Reset()
			} else if sig == syscall.Signal(0x1e /**syscall.SIGUSR1**/) { // 刷新代理状态
				err := teaproxy.SharedManager.Restart()
				if err != nil {
					logs.Println("[error]" + err.Error())
				} else {
					proxyutils.FinishChange()
				}
			} else if sig == syscall.Signal(0x1f /**syscall.SIGUSR2**/) { // 同步
				node := teaconfigs.SharedNodeConfig()
				if node == nil {
					logs.Println("[cluster]not a node yet")
					return
				}

				if node.IsMaster() {
					logs.Println("[cluster]push items")
					teacluster.SharedManager.BuildSum()
					teacluster.SharedManager.PushItems()
				} else {
					logs.Println("[cluster]pull items")
					teacluster.SharedManager.BuildSum()
					teacluster.SharedManager.PullItems()
				}
			} else {
				if sig == syscall.SIGINT { // 终止进程
					if server != nil {
						server.Stop()
						time.Sleep(1 * time.Second)
					}
				}

				// 删除PID
				pidFile := files.NewFile(Tea.Root + "/bin/pid")
				if pidFile.Exists() {
					err = pidFile.Delete()
					if err != nil {
						logs.Error(err)
					}
				}
				os.Exit(0)
			}
		}
	}()
}

// 重置Root
func (this *WebShell) resetRoot() {
	if !Tea.IsTesting() {
		exePath, err := os.Executable()
		if err != nil {
			exePath = os.Args[0]
		}
		link, err := filepath.EvalSymlinks(exePath)
		if err == nil {
			exePath = link
		}
		fullPath, err := filepath.Abs(exePath)
		if err == nil {
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(fullPath)))
		}
	}
	Tea.SetPublicDir(Tea.Root + Tea.DS + "web" + Tea.DS + "public")
	Tea.SetViewsDir(Tea.Root + Tea.DS + "web" + Tea.DS + "views")
	Tea.SetTmpDir(Tea.Root + Tea.DS + "web" + Tea.DS + "tmp")
}

// 检查命令行参数
func (this *WebShell) execArgs(writer io.Writer) bool {
	if len(os.Args) == 1 {
		// 检查是否已经启动
		proc := this.checkPid()
		if proc != nil {
			this.write(writer, "TeaWeb is already running, pid:", proc.Pid)
			return true
		}
		return false
	}
	args := os.Args[1:]
	arg0 := ""
	if len(args) > 0 {
		arg0 = args[0]
	}
	if this.hasArg(arg0, "?", "help", "-help", "h", "-h") { // 帮助
		return this.ExecHelp(writer)
	} else if this.hasArg(arg0, "-v", "version", "-version") { // 版本号
		return this.ExecVersion(writer)
	} else if this.hasArg(arg0, "start") { // 启动
		return this.ExecStart(writer)
	} else if this.hasArg(arg0, "stop") { // 停止
		return this.ExecStop(os.Stdout)
	} else if this.hasArg(arg0, "reload") { // 重新加载代理配置
		return this.ExecReload(writer)
	} else if this.hasArg(arg0, "restart") { // 重启
		return this.ExecRestart(writer)
	} else if this.hasArg(arg0, "reset") { // 重置
		return this.ExecReset(writer)
	} else if this.hasArg(arg0, "status") { // 状态
		return this.ExecStatus(writer)
	} else if this.hasArg(arg0, "sync") { // 同步
		return this.ExecSync(writer)
	} else if this.hasArg(arg0, "service") && runtime.GOOS == "windows" { // Windows服务
		return this.ExecService(writer)
	} else if this.hasArg(arg0, "pprof") {
		return this.ExecPprof(writer)
	}

	if len(args) > 0 {
		this.write(writer, "Unknown command option '"+strings.Join(args, " ")+"', run './bin/teaweb -h' to lookup the usage.")
		return true
	}
	return false
}

// 帮助
func (this *WebShell) ExecHelp(writer io.Writer) bool {
	this.write(writer, "TeaWeb v"+teaconst.TeaVersion)
	this.write(writer, "Usage:", "\n   ./bin/teaweb [option]")
	this.write(writer, "")
	this.write(writer, "Options:")
	this.write(writer, "  -h", "\n     print this help")
	this.write(writer, "  -v", "\n     print version")
	this.write(writer, "  start", "\n     start the server in background")
	this.write(writer, "  stop", "\n     stop the server")
	this.write(writer, "  reload", "\n     reload all proxy servers configs")
	this.write(writer, "  restart", "\n     restart the server")
	this.write(writer, "  reset", "\n     reset the server locker status")
	this.write(writer, "  status", "\n     print server status")
	this.write(writer, "  sync", "\n     sync config files with cluster")
	this.write(writer, "  pprof [address]", "\n     start pprof server")
	this.write(writer, "")
	this.write(writer, "To run the server in foreground:", "\n   ./bin/teaweb")

	return true
}

// 版本号
func (this *WebShell) ExecVersion(writer io.Writer) bool {
	this.write(writer, "TeaWeb v"+teaconst.TeaVersion, "(build: "+runtime.Version(), runtime.GOOS, runtime.GOARCH+")")
	return true
}

// 启动
func (this *WebShell) ExecStart(writer io.Writer) bool {
	proc := this.checkPid()
	if proc != nil {
		this.write(writer, "TeaWeb already started, pid:", proc.Pid)
		return true
	}

	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		this.write(writer, "TeaWeb  start failed:", err.Error())
		return true
	}
	this.write(writer, "TeaWeb started ok, pid:", cmd.Process.Pid)

	return true
}

// 停止
func (this *WebShell) ExecStop(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, "TeaWeb not started")
		return true
	}

	err := proc.Kill()
	if err != nil {
		this.write(writer, "TeaWeb stop error:", err.Error())
		return true
	}

	err = files.NewFile(Tea.Root + "/bin/pid").Delete()
	this.write(writer, "TeaWeb stopped ok, pid:", proc.Pid)

	if err != nil {
		this.write(writer, "ERROR:", err.Error())
	}

	return true
}

// 重载代理配置
func (this *WebShell) ExecReload(writer io.Writer) bool {
	pidString, err := files.NewFile(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid").ReadAllString()
	if err != nil {
		this.write(writer, err.Error())
		return true
	}

	pid := types.Int(pidString)
	proc, err := os.FindProcess(pid)
	if err != nil {
		this.write(writer, err.Error())
		return true
	}
	if proc == nil {
		this.write(writer, "can not find process")
		return true
	}
	err = proc.Signal(syscall.Signal(0x1e /**syscall.SIGUSR1**/))
	if err != nil {
		logs.Error(err)
		return true
	}
	this.write(writer, "reload success")
	return true
}

// 重启
func (this *WebShell) ExecRestart(writer io.Writer) bool {
	proc := this.checkPid()
	if proc != nil {
		err := proc.Kill()
		if err != nil {
			this.write(writer, "TeaWeb stop error:", err.Error())
			return true
		}

		// 等待进程结束
		time.Sleep(1 * time.Second)
	}

	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		this.write(writer, "TeaWeb restart failed:", err.Error())
		return true
	}
	this.write(writer, "TeaWeb restarted ok, pid:", cmd.Process.Pid)

	return true
}

// 重置
func (this *WebShell) ExecReset(writer io.Writer) bool {
	pidString, err := files.NewFile(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid").ReadAllString()
	if err != nil {
		this.write(writer, err.Error())
		return true
	}

	pid := types.Int(pidString)
	proc, err := os.FindProcess(pid)
	if err != nil {
		this.write(writer, err.Error())
		return true
	}
	if proc == nil {
		this.write(writer, "can not find process")
		return true
	}
	err = proc.Signal(syscall.SIGHUP)
	if err != nil {
		this.write(writer, err.Error())
		return true
	}
	this.write(writer, "reset success")
	return true
}

// 状态
func (this *WebShell) ExecStatus(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, "TeaWeb not started yet")
	} else {
		this.write(writer, "TeaWeb is running, pid:"+fmt.Sprintf("%d", proc.Pid))
	}
	return true
}

// 同步
func (this *WebShell) ExecSync(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, "TeaWeb not started yet")
	} else {
		err := proc.Signal(syscall.Signal(0x1f /**syscall.SIGUSR2**/))
		if err != nil {
			logs.Error(err)
			return true
		}
		this.write(writer, "signal sent successfully")
	}
	return true
}

// 启动pprof
func (this *WebShell) ExecPprof(writer io.Writer) bool {
	addr := ":6060"
	if len(os.Args) == 3 {
		addr = os.Args[2]
	}
	this.write(writer, "===\nstart pprof server '"+addr+"'\n===")
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			this.write(writer, "[error]"+err.Error())
		}
	}()

	return false
}

// 检查PID
func (this *WebShell) checkPid() *os.Process {
	// check pid file
	pidFile := files.NewFile(Tea.Root + "/bin/pid")
	if !pidFile.Exists() {
		return nil
	}
	pidString, err := pidFile.ReadAllString()
	if err != nil {
		return nil
	}
	pid := types.Int(pidString)

	if pid <= 0 {
		return nil
	}

	// 如果是当前进程在检查，说明没有启动
	if pid == os.Getpid() {
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil || proc == nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		return proc
	}

	err = proc.Signal(syscall.Signal(0)) // 根据方法文档：Sending Interrupt on Windows is not implemented
	if err != nil {
		return nil
	}

	// ps?
	ps, err := exec.LookPath("ps")
	if err != nil {
		return proc
	}

	cmd := exec.Command(ps, "-p", pidString, "-o", "command=")
	output, err := cmd.Output()
	if err != nil {
		return proc
	}

	if len(output) == 0 {
		return nil
	}

	outputString := string(output)
	index := strings.LastIndex(outputString, "/")
	if index > -1 {
		outputString = outputString[index+1:]
	}
	index2 := strings.LastIndex(outputString, "\\")
	if index2 > 0 {
		outputString = outputString[index2+1:]
	}
	if strings.Contains(outputString, "teaweb") && !strings.Contains(outputString, "teaweb-") {
		return proc
	}

	return nil
}

// 写入string到writer
func (this *WebShell) write(writer io.Writer, args ...interface{}) {
	_, _ = fmt.Fprintln(writer, args...)
}

// 判断命令
func (this *WebShell) hasArg(arg string, value ...string) bool {
	return lists.ContainsAny(value, arg)
}
