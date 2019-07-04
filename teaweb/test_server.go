package teaweb

import (
	"bytes"
	"fmt"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/TeaGo/utils/time"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

// 测试服务器
func startTestServer() {
	TeaGo.NewServer(false).
		AccessLog(false).
		Get("/", func(resp http.ResponseWriter) {
			resp.Write([]byte("This is test server"))
		}).
		Get("/hello", func(req *http.Request, resp http.ResponseWriter) {
			resp.Write([]byte(req.RequestURI + ":"))
			resp.Write([]byte("world"))
		}).
		Get("/benchmark", func(resp http.ResponseWriter) {
			resp.Write([]byte("Hello, World, this is benchmark url"))
		}).
		Get("/redirect", func(req *http.Request, resp http.ResponseWriter) {
			code := types.Int(req.URL.Query().Get("code"))
			if code >= 300 && code < 400 {
				resp.Header().Set("Location", "/redirect2")
				resp.Header().Set("Set-Cookie", "code="+fmt.Sprintf("%d", code)+"; Max-Agent=86400; Path=/")
				resp.WriteHeader(code)
			} else {
				http.Redirect(resp, req, "/redirect2", http.StatusTemporaryRedirect)
			}
		}).
		Get("/redirect2", func(req *http.Request, resp http.ResponseWriter) {
			for k, v := range req.Header {
				for _, v1 := range v {
					resp.Write([]byte( k + ": " + v1 + "\n"))
				}
			}

			resp.Write([]byte("\n\n"))
			resp.Write([]byte("the page after redirect"))
		}).
		Get("/webhook", func(req *http.Request, resp http.ResponseWriter) {
			resp.Write([]byte("Get " + req.URL.String() + "\n"))
			for k, v := range req.Header {
				for _, v1 := range v {
					resp.Write([]byte( k + ": " + v1 + "\n"))
				}
			}

			// 测试超过1560内容长度
			resp.Write(bytes.Repeat([]byte{' '}, 934))
			resp.Write([]byte{'a'})
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
		Post("/upload", func(req *http.Request, resp http.ResponseWriter) {
			err := req.ParseMultipartForm(32 * 1024 * 1024)
			if err != nil {
				resp.Write([]byte(err.Error()))
				return
			}

			resp.Write([]byte("files:\n"))
			for field, files := range req.MultipartForm.File {
				for _, f := range files {
					resp.Write([]byte(field + ":" + f.Filename + ", " + fmt.Sprintf("%d", f.Size) + "bytes\n"))
				}
			}

			resp.Write([]byte("params:\n"))
			for k, values := range req.PostForm {
				for _, v := range values {
					resp.Write([]byte(k + ":" + v + "\n"))
				}
			}

		}).
		Get("/cookie", func(req *http.Request, resp http.ResponseWriter) {
			resp.Header().Add("Set-Cookie", "Asset_UserId=1; expires=Sun, 05-May-2019 14:42:21 GMT; path=/", )
			resp.Write([]byte("set cookie"))
		}).
		GetPost("/json", func(req *http.Request, resp http.ResponseWriter) {
			resp.Header().Set("Content-Type", "application/json")
			data := maps.Map{
				"hello": "world",
			}
			resp.Write([]byte(stringutil.JSONEncode(data)))
		}).
		Get("/nocache", func(req *http.Request, resp http.ResponseWriter) {
			resp.Header().Set("Cache-Control", "no-cache")
			resp.Write([]byte("will be not cached " + timeutil.Format("Y-m-d H:i:s")))
		}).
		Get("/gzip", func(req *http.Request, resp http.ResponseWriter) {
			compressResource(resp, Tea.PublicDir()+"/js/vue.min.js", "text/javascript; charset=utf-8")
		}).
		Get("/image", func(req *http.Request, resp http.ResponseWriter) {
			data, err := ioutil.ReadFile(Tea.PublicDir() + "/images/logo.png")
			if err != nil {
				resp.Write([]byte(err.Error()))
			} else {
				resp.Header().Set("Content-Type", "image/png")
				resp.Write(data)
			}
		}).
		Put("/put", func(req *http.Request, resp http.ResponseWriter) {
			data, err := httputil.DumpRequest(req, true)
			if err != nil {
				resp.Write([]byte(err.Error()))
				return
			}
			resp.Write(data)
		}).
		StartOn("127.0.0.1:9991")
}
