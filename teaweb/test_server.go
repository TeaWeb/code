package teaweb

import (
	"github.com/iwind/TeaGo"
	"net/http"
)

// 测试服务器
func startTestServer() {
	TeaGo.NewServer(false).
		AccessLog(false).
		Get("/", func(resp http.ResponseWriter) {
			resp.Write([]byte("This is test server"))
		}).
		Get("/benchmark", func(resp http.ResponseWriter) {
			resp.Write([]byte("Hello, World, this is benchmark url"))
		}).
		Get("/redirect", func(req *http.Request, resp http.ResponseWriter) {
			http.Redirect(resp, req, "/redirect2", http.StatusTemporaryRedirect)
		}).
		Get("/redirect2", func(resp http.ResponseWriter) {
			resp.Write([]byte("the page after redirect"))
		}).
		StartOn("127.0.0.1:9991")
}
