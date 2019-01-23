package teaweb

import (
	"github.com/iwind/TeaGo"
	"io/ioutil"
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
		Get("/webhook", func(req *http.Request, resp http.ResponseWriter) {
			resp.Write([]byte("call " + req.URL.String()))
		}).
		Post("/webhook", func(req *http.Request, resp http.ResponseWriter) {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				resp.Write([]byte("error:" + err.Error()))
			} else {
				resp.Write([]byte("post:" + string(body)))
			}
		}).
		StartOn("127.0.0.1:9991")
}
