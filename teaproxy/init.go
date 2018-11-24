package teaproxy

import "github.com/TeaWeb/code/teaconfigs"

// 所有监听器集合
var LISTENERS = []*Listener{}

// 所有服务
var SERVERS = map[string]*teaconfigs.ServerConfig{} // id => server
