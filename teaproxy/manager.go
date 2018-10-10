package teaproxy

import (
	"sync"
	"github.com/iwind/TeaGo/logs"
	"github.com/TeaWeb/code/teaconfigs"
	_ "github.com/TeaWeb/code/teastats" // 引入统计处理工具
)

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

func Wait() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func Shutdown() {
	for _, listener := range LISTENERS {
		listener.Shutdown()
	}

	LISTENERS = []*Listener{}
	SERVERS = map[string]*teaconfigs.ServerConfig{}
}

func Restart() {
	Shutdown()
	Start()
}

func FindServer(id string) (server *teaconfigs.ServerConfig, found bool) {
	server, found = SERVERS[id]
	return
}
