package teaproxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"net/http"
)

// 所有监听器集合
var LISTENERS = []*Listener{}

// 所有服务
var SERVERS = map[string]*teaconfigs.ServerConfig{} // id => server

// 状态码筛选
var StatusCodeParser func(statusCode int, headers http.Header, respData []byte, parserScript string) (string, error) = nil
