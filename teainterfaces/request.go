package teainterfaces

import "net/http"

type PluginRequestFilterInterface interface {
	FilterRequest(request *http.Request) bool // 过滤请求，如果返回false，则不会往下执行
}
