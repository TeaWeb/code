package teaweb

import (
	"github.com/iwind/TeaGo"
	"io/ioutil"
	"net/http"
	"time"
)

// 测试服务器
func startTestServer() {
	TeaGo.NewServer(false).
		AccessLog(false).
		Get("/", func(resp http.ResponseWriter) {
			resp.Write([]byte("This is test server"))
		}).
		Get("/hello", func(resp http.ResponseWriter) {
			resp.Write([]byte("world"))
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
			for k, v := range req.Header {
				for _, v1 := range v {
					resp.Write([]byte("Header " + k + " " + v1 + "\n"))
				}
			}
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				resp.Write([]byte("error:" + err.Error()))
			} else {
				resp.Write([]byte("post:" + string(body)))
			}
		}).
		Get("/timeout30", func(req *http.Request, resp http.ResponseWriter) {
			time.Sleep(31 * time.Second)
			resp.Write([]byte("30 seconds timeout"))
		}).
		Get("/timeout120", func(req *http.Request, resp http.ResponseWriter) {
			time.Sleep(121 * time.Second)
			resp.Write([]byte("120 seconds timeout"))
		}).
		StartOn("127.0.0.1:9991")
}
