package teainterfaces

import "net/http"

type PluginResponseFilterInterface interface {
	FilterResponse(response *http.Response, writer http.ResponseWriter) bool // 过滤响应，如果返回false，则不会往下执行
}
