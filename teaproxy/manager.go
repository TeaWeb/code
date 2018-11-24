package teaproxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	_ "github.com/TeaWeb/code/teastats" // 引入统计处理工具
	"github.com/iwind/TeaGo/logs"
	"sync"
)

// 启动服务
func Start() {
	listenerConfigs, err := teaconfigs.ParseConfigs()
	if err != nil {
		logs.Error(err)
		return
	}

	for _, config := range listenerConfigs {
		for _, s := range config.Servers {
			SERVERS[s.Id] = s
		}

		listener := NewListener(config)
		go listener.Start()
	}
}

// 等待服务执行完毕
func Wait() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

// 关闭服务
func Shutdown() {
	for _, listener := range LISTENERS {
		listener.Shutdown()
	}

	LISTENERS = []*Listener{}
	SERVERS = map[string]*teaconfigs.ServerConfig{}
}

// 重启服务
func Restart() {
	Shutdown()
	Start()
}

// 查找服务
func FindServer(id string) (server *teaconfigs.ServerConfig, found bool) {
	server, found = SERVERS[id]
	return
}
