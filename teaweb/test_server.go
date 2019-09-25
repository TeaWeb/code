package teaweb

import (
	"fmt"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
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
			_, _ = resp.Write([]byte("This is test server"))
		}).
		Get("/hello", func(req *http.Request, resp http.ResponseWriter) {
			_, _ = resp.Write([]byte(req.RequestURI + ":"))
			_, _ = resp.Write([]byte("world"))
		}).
		Get("/benchmark", func(resp http.ResponseWriter) {
			_, _ = resp.Write([]byte("Hello, World, this is benchmark url"))
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
					_, _ = resp.Write([]byte( k + ": " + v1 + "\n"))
				}
			}

			_, _ = resp.Write([]byte("\n\n"))
			_, _ = resp.Write([]byte("the page after redirect"))
		}).
		Get("/webhook", func(req *http.Request, resp http.ResponseWriter) {
			_, _ = resp.Write([]byte("Get " + req.URL.String() + "\n"))
			for k, v := range req.Header {
				for _, v1 := range v {
					_, _ = resp.Write([]byte( k + ": " + v1 + "\n"))
				}
			}
		}).
		Post("/webhook", func(req *http.Request, resp http.ResponseWriter) {
			for k, v := range req.Header {
				for _, v1 := range v {
					_, _ = resp.Write([]byte("Header " + k + ": " + v1 + "\n"))
				}
			}
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				_, _ = resp.Write([]byte("error: " + err.Error()))
			} else {
				_, _ = resp.Write([]byte("post: " + string(body)))
			}
		}).
		Put("/webhook", func(req *http.Request, resp http.ResponseWriter) {
			for k, v := range req.Header {
				for _, v1 := range v {
					_, err := resp.Write([]byte("Header " + k + ": " + v1 + "\n"))
					if err != nil {
						logs.Error(err)
					}
				}
			}
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				_, err = resp.Write([]byte("error: " + err.Error()))
			} else {
				_, err = resp.Write([]byte("put: " + string(body)))
			}
		}).
		Get("/timeout30", func(req *http.Request, resp http.ResponseWriter) {
			time.Sleep(31 * time.Second)
			_, _ = resp.Write([]byte("30 seconds timeout"))
		}).
		Get("/timeout120", func(req *http.Request, resp http.ResponseWriter) {
			time.Sleep(121 * time.Second)
			_, _ = resp.Write([]byte("120 seconds timeout"))
		}).
		Post("/upload", func(req *http.Request, resp http.ResponseWriter) {
			err := req.ParseMultipartForm(32 * 1024 * 1024)
			if err != nil {
				_, _ = resp.Write([]byte(err.Error()))
				return
			}

			_, _ = resp.Write([]byte("files:\n"))
			for field, files := range req.MultipartForm.File {
				for _, f := range files {
					_, _ = resp.Write([]byte(field + ":" + f.Filename + ", " + fmt.Sprintf("%d", f.Size) + "bytes\n"))
				}
			}

			_, _ = resp.Write([]byte("params:\n"))
			for k, values := range req.PostForm {
				for _, v := range values {
					_, _ = resp.Write([]byte(k + ":" + v + "\n"))
				}
			}

		}).
		Get("/cookie", func(req *http.Request, resp http.ResponseWriter) {
			resp.Header().Add("Set-Cookie", "Asset_UserId=1; expires=Sun, 05-May-2019 14:42:21 GMT; path=/", )
			_, _ = resp.Write([]byte("set cookie"))
		}).
		GetPost("/json", func(req *http.Request, resp http.ResponseWriter) {
			resp.Header().Set("Content-Type", "application/json")
			data := maps.Map{
				"hello": "world",
			}
			_, _ = resp.Write([]byte(stringutil.JSONEncode(data)))
		}).
		Get("/nocache", func(req *http.Request, resp http.ResponseWriter) {
			resp.Header().Set("Cache-Control", "no-cache")
			_, _ = resp.Write([]byte("will be not cached " + timeutil.Format("Y-m-d H:i:s")))
		}).
		Get("/gzip", func(req *http.Request, resp http.ResponseWriter) {
			compressResource(resp, Tea.PublicDir()+"/js/vue.min.js", "text/javascript; charset=utf-8")
		}).
		Get("/image", func(req *http.Request, resp http.ResponseWriter) {
			data, err := ioutil.ReadFile(Tea.PublicDir() + "/images/logo.png")
			if err != nil {
				_, _ = resp.Write([]byte(err.Error()))
			} else {
				resp.Header().Set("Content-Type", "image/png")
				_, _ = resp.Write(data)
			}
		}).
		Put("/put", func(req *http.Request, resp http.ResponseWriter) {
			data, err := httputil.DumpRequest(req, true)
			if err != nil {
				_, _ = resp.Write([]byte(err.Error()))
				return
			}
			_, _ = resp.Write(data)
		}).
		Get("/basicAuth", func(req *http.Request, resp http.ResponseWriter) {
			if len(req.Header.Get("Authorization")) == 0 {
				resp.Header().Set("WWW-Authenticate", `Basic realm="My Realm"`)
				resp.WriteHeader(401)
			}

			for k, v := range req.Header {
				for _, v1 := range v {
					_, _ = resp.Write([]byte("Header " + k + " " + v1 + "\n"))
				}
			}
		}).
		Get("/html", func(req *http.Request, resp http.ResponseWriter) {
			_, _ = resp.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>HTML Page</title>
</head>
<body>
<strong>THIS IS HTML BODY</strong>
</body>
</html>
`))
		}).
		StartOn("127.0.0.1:9991")
}
