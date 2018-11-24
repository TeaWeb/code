package teaproxy

import "github.com/TeaWeb/code/teaconfigs"

// 开启API功能
func EnableAPIServer(serverId string) {
	server, found := SERVERS[serverId]
	if found {
		server.APIOn = true
	}
}

// 添加API
func AddAPI(serverId string, api *teaconfigs.API) {
	server, found := SERVERS[serverId]
	if found {
		server.AddAPI(api)
	}
}

// 删除API
func DeleteAPI(serverId string, api *teaconfigs.API) {
	server, found := SERVERS[serverId]
	if found {
		server.DeleteAPI(api)
	}
}
