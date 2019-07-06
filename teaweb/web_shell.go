package teaweb

import (
	"fmt"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// 命令行相关封装
type WebShell struct {
	ShouldStop bool
}

// 启动
func (this *WebShell) Start() {
	// 重置ROOT
	this.resetRoot()

	// 执行参数
	if this.execArgs() {
		this.ShouldStop = true
		return
	}

	// 当前PID
	files.NewFile(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid").
		WriteString(fmt.Sprintf("%d", os.Getpid()))

	// 信号
	signalsChannel := make(chan os.Signal, 1024)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGHUP, syscall.Signal(0x1e) /**syscall.SIGUSR1**/, syscall.SIGTERM)
	go func() {
		for {
			sig := <-signalsChannel

			if sig == syscall.SIGHUP { // 重置
				configs.SharedAdminConfig().Reset()
			} else if sig == syscall.Signal(0x1e) /**syscall.SIGUSR1**/ { // 刷新代理状态
				err := teaproxy.SharedManager.Restart()
				if err != nil {
					logs.Println("[error]" + err.Error())
				} else {
					proxyutils.FinishChange()
				}
			} else {
				if sig == syscall.SIGINT {
					if server != nil {
						server.Stop()
						time.Sleep(1 * time.Second)
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
func (this *WebShell) execArgs() bool {
	if len(os.Args) == 1 {
		// 检查是否已经启动
		proc := this.checkPid()
		if proc != nil {
			fmt.Println("TeaWeb is already running, pid:", proc.Pid)
			return true
		}
		return false
	}
	args := os.Args[1:]
	if lists.ContainsAny(args, "?", "help", "-help", "h", "-h") { // 帮助
		return this.execHelp()
	} else if lists.ContainsAny(args, "-v", "version", "-version") { // 版本号
		return this.execVersion()
	} else if lists.ContainsString(args, "start") { // 启动
		return this.execStart()
	} else if lists.ContainsString(args, "stop") { // 停止
		return this.execStop()
	} else if lists.ContainsString(args, "reload") { // 重新加载代理配置
		return this.execReload()
	} else if lists.ContainsString(args, "restart") { // 重启
		return this.execRestart()
	} else if lists.ContainsString(args, "reset") { // 重置
		return this.execReset()
	} else if lists.ContainsString(args, "status") { // 状态
		return this.execStatus()
	} else if lists.ContainsString(args, "service") && runtime.GOOS == "windows" { // Windows服务
		return this.execService()
	}

	if len(args) > 0 {
		fmt.Println("Unknown command option '" + strings.Join(args, " ") + "', run './bin/teaweb -h' to lookup the usage.")
		return true
	}
	return false
}

// 帮助
func (this *WebShell) execHelp() bool {
	fmt.Println("TeaWeb v" + teaconst.TeaVersion)
	fmt.Println("Usage:", "\n   ./bin/teaweb [option]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h", "\n     print this help")
	fmt.Println("  -v", "\n     print version")
	fmt.Println("  start", "\n     start the server in background")
	fmt.Println("  stop", "\n     stop the server")
	fmt.Println("  reload", "\n     reload all proxy servers configs")
	fmt.Println("  restart", "\n     restart the server")
	fmt.Println("  reset", "\n     reset the server locker status")
	fmt.Println("  status", "\n     print server status")
	fmt.Println("")
	fmt.Println("To run the server in foreground:", "\n   ./bin/teaweb")

	return true
}

// 版本号
func (this *WebShell) execVersion() bool {
	fmt.Println("TeaWeb v"+teaconst.TeaVersion, "(build: "+runtime.Version(), runtime.GOOS, runtime.GOARCH+")")
	return true
}

// 启动
func (this *WebShell) execStart() bool {
	proc := this.checkPid()
	if proc != nil {
		fmt.Println("TeaWeb already started, pid:", proc.Pid)
		return true
	}

	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		fmt.Println("TeaWeb  start failed:", err.Error())
		return true
	}
	fmt.Println("TeaWeb started ok, pid:", cmd.Process.Pid)

	return true
}

// 停止
func (this *WebShell) execStop() bool {
	proc := this.checkPid()
	if proc == nil {
		fmt.Println("TeaWeb not started")
		return true
	}

	err := proc.Kill()
	if err != nil {
		fmt.Println("TeaWeb stop error:", err.Error())
		return true
	}

	files.NewFile(Tea.Root + "/bin/pid").Delete()
	fmt.Println("TeaWeb stopped ok, pid:", proc.Pid)

	return true
}

// 重载代理配置
func (this *WebShell) execReload() bool {
	pidString, err := files.NewFile(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid").ReadAllString()
	if err != nil {
		logs.Error(err)
		return true
	}

	pid := types.Int(pidString)
	proc, err := os.FindProcess(pid)
	if err != nil {
		logs.Error(err)
		return true
	}
	if proc == nil {
		logs.Println("can not find process")
		return true
	}
	err = proc.Signal(syscall.Signal(0x1e) /**syscall.SIGUSR1**/)
	if err != nil {
		logs.Error(err)
		return true
	}
	fmt.Println("reload success")
	return true
}

// 重启
func (this *WebShell) execRestart() bool {
	proc := this.checkPid()
	if proc != nil {
		err := proc.Kill()
		if err != nil {
			fmt.Println("TeaWeb stop error:", err.Error())
			return true
		}
	}

	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		fmt.Println("TeaWeb restart failed:", err.Error())
		return true
	}
	fmt.Println("TeaWeb restarted ok, pid:", cmd.Process.Pid)

	return true
}

// 重置
func (this *WebShell) execReset() bool {
	pidString, err := files.NewFile(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid").ReadAllString()
	if err != nil {
		logs.Error(err)
		return true
	}

	pid := types.Int(pidString)
	proc, err := os.FindProcess(pid)
	if err != nil {
		logs.Error(err)
		return true
	}
	if proc == nil {
		logs.Println("can not find process")
		return true
	}
	err = proc.Signal(syscall.SIGHUP)
	if err != nil {
		logs.Error(err)
		return true
	}
	fmt.Println("reset success")
	return true
}

// 状态
func (this *WebShell) execStatus() bool {
	proc := this.checkPid()
	if proc == nil {
		fmt.Println("TeaWeb not started yet")
	} else {
		fmt.Println("TeaWeb is running, pid:" + fmt.Sprintf("%d", proc.Pid))
	}
	return true
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
