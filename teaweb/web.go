package teaweb

import (
	"fmt"
	_ "github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teaproxy"
	_ "github.com/TeaWeb/code/teaweb/actions/default/agents"
	_ "github.com/TeaWeb/code/teaweb/actions/default/agents/apps"
	_ "github.com/TeaWeb/code/teaweb/actions/default/agents/board"
	_ "github.com/TeaWeb/code/teaweb/actions/default/agents/notices"
	_ "github.com/TeaWeb/code/teaweb/actions/default/agents/settings"
	_ "github.com/TeaWeb/code/teaweb/actions/default/api/agent"
	_ "github.com/TeaWeb/code/teaweb/actions/default/api/monitor"
	_ "github.com/TeaWeb/code/teaweb/actions/default/cache"
	_ "github.com/TeaWeb/code/teaweb/actions/default/dashboard"
	"github.com/TeaWeb/code/teaweb/actions/default/index"
	_ "github.com/TeaWeb/code/teaweb/actions/default/log"
	_ "github.com/TeaWeb/code/teaweb/actions/default/login"
	"github.com/TeaWeb/code/teaweb/actions/default/logout"
	_ "github.com/TeaWeb/code/teaweb/actions/default/mongo"
	_ "github.com/TeaWeb/code/teaweb/actions/default/plugins"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/backend"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/board"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/fastcgi"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/headers"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/locations"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/websocket"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/log"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/rewrite"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/ssl"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/login"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/mongo"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/profile"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/server"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/update"
	_ "github.com/TeaWeb/code/teaweb/actions/default/stat"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/sessions"
	"github.com/iwind/TeaGo/types"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// 启动
func Start() {
	if lookupArgs() {
		return
	}

	// 信号
	signalsChannel := make(chan os.Signal, 1024)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		for {
			sig := <-signalsChannel

			// 重置
			if sig == syscall.SIGHUP {
				configs.SharedAdminConfig().Reset()
			} else {
				os.Exit(0)
			}
		}
	}()

	// 日志
	writer := new(utils.LogWriter)
	writer.Init()
	logs.SetWriter(writer)

	// 启动代理
	go func() {
		time.Sleep(1 * time.Second)

		// 启动代理
		teaproxy.Start()
	}()

	// 启动测试服务器
	if Tea.IsTesting() {
		go func() {
			time.Sleep(1 * time.Second)

			//startTestServer()
		}()
	}

	// 启动管理界面
	TeaGo.NewServer().
		AccessLog(false).

		Get("/", new(index.IndexAction)).
		Get("/logout", new(logout.IndexAction)).

		EndAll().

		Session(sessions.NewFileSessionManager(
			86400,
			"gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9",
		)).

		Start()
}

// 检查命令行参数
func lookupArgs() bool {
	if len(os.Args) == 1 {
		return false
	}
	args := os.Args[1:]
	if lists.ContainsAny(args, "?", "help", "-help", "h", "-h") { // 帮助
		fmt.Println("TeaWeb v" + teaconst.TeaVersion)
		fmt.Println("Usage:", "\n   ./bin/teaweb [option]")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  -h", "\n     print this help")
		fmt.Println("  -v", "\n     print version")
		fmt.Println("  start", "\n     start the server")
		fmt.Println("  stop", "\n     stop the server")
		fmt.Println("  restart", "\n     restart the server")
		fmt.Println("  reset", "\n     reset the server status")
		return true
	} else if lists.Contains(args, "-v") { // 版本号
		fmt.Println("TeaWeb v"+teaconst.TeaVersion, "(build: "+runtime.Version(), runtime.GOOS, runtime.GOARCH+")")
		return true
	} else if lists.Contains(args, "start") { // 启动
		proc := checkPid()
		if proc != nil {
			fmt.Println("[teaweb]already started, pid:", proc.Pid)
			return true
		}

		cmd := exec.Command(os.Args[0])
		err := cmd.Start()
		if err != nil {
			fmt.Println("[teaweb]start failed:", err.Error())
			return true
		}
		fmt.Println("[teaweb]started ok, pid:", cmd.Process.Pid)

		return true
	} else if lists.Contains(args, "stop") { // 停止
		proc := checkPid()
		if proc == nil {
			fmt.Println("[teaweb]not started")
			return true
		}

		err := proc.Kill()
		if err != nil {
			fmt.Println("[teaweb]stop error:", err.Error())
			return true
		}

		files.NewFile(Tea.Root + "/bin/pid").Delete()
		fmt.Println("[teaweb]stopped ok, pid:", proc.Pid)

		return true
	} else if lists.Contains(args, "restart") { // 重启
		proc := checkPid()
		if proc != nil {
			err := proc.Kill()
			if err != nil {
				fmt.Println("[teaweb]stop error:", err.Error())
				return true
			}
		}

		cmd := exec.Command(os.Args[0])
		err := cmd.Start()
		if err != nil {
			fmt.Println("[teaweb]restart failed:", err.Error())
			return true
		}
		fmt.Println("[teaweb]restarted ok, pid:", cmd.Process.Pid)

		return true
	} else if lists.Contains(args, "reset") { // 重置
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
		logs.Println("reset success")
		return true
	}

	if len(args) > 0 {
		fmt.Println("[teaweb]unknown command option '" + strings.Join(args, " ") + "', run './bin/teaweb -h' to see the usage.")
		return true
	}
	return false
}

// 检查PID
func checkPid() *os.Process {
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
	proc, err := os.FindProcess(pid)
	if err != nil || proc == nil {
		return nil
	}

	err = proc.Signal(syscall.Signal(0))
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

	if strings.Contains(string(output), "teaweb") {
		return proc
	}

	return nil
}
