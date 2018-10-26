package teaproxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/gofcgi"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var requestVarReg = regexp.MustCompile("\\${[\\w.-]+}")

// 文本mime-type列表
var textMimeMap = map[string]bool{
	"application/atom+xml":                true,
	"application/javascript":              true,
	"application/x-javascript":            true,
	"application/json":                    true,
	"application/rss+xml":                 true,
	"application/x-web-app-manifest+json": true,
	"application/xhtml+xml":               true,
	"application/xml":                     true,
	"image/svg+xml":                       true,
	"text/css":                            true,
	"text/plain":                          true,
	"text/javascript":                     true,
	"text/xml":                            true,
	"text/html":                           true,
	"text/xhtml":                          true,
	"text/sgml":                           true,
}

// 请求定义
type Request struct {
	raw    *http.Request
	server *teaconfigs.ServerConfig

	scheme        string
	rawScheme     string // 原始的scheme
	uri           string
	rawURI        string // 跳转之前的uri
	host          string
	method        string
	serverName    string // @TODO
	serverAddr    string
	charset       string
	headers       []*teaconfigs.HeaderConfig // 自定义Header
	ignoreHeaders []string                   // 忽略的Header
	varMapping    map[string]string          // 自定义变量

	root     string   // 资源根目录
	index    []string // 目录下默认访问的文件
	backend  *teaconfigs.ServerBackendConfig
	fastcgi  *teaconfigs.FastcgiConfig
	proxy    *teaconfigs.ServerConfig
	location *teaconfigs.LocationConfig

	rewriteId      string // 匹配的rewrite id
	rewriteReplace string // 经过rewrite之后的URL

	// 执行请求
	filePath string

	responseBytesSent     int64  // @TODO
	responseBodyBytesSent int64  // @TODO
	responseStatus        int    // @TODO
	responseStatusMessage string // @TODO
	responseHeader        map[string][]string

	requestFromTime    time.Time
	requestTime        float64 // @TODO
	requestTimeISO8601 string
	requestTimeLocal   string
	requestMsec        float64
	requestTimestamp   int64

	shouldLog bool
	debug     bool
}

// 获取新的请求
func NewRequest(rawRequest *http.Request) *Request {
	now := time.Now()
	return &Request{
		varMapping:         map[string]string{},
		raw:                rawRequest,
		rawURI:             rawRequest.URL.RequestURI(),
		requestFromTime:    now,
		requestTimestamp:   now.Unix(),
		requestTimeISO8601: now.Format("2006-01-02T15:04:05.000Z07:00"),
		requestTimeLocal:   now.Format("2/Jan/2006:15:04:05 -0700"),
		requestMsec:        float64(now.Unix()) + float64(now.Nanosecond())/1000000000,
		shouldLog:          true,
	}
}

func (this *Request) configure(server *teaconfigs.ServerConfig, redirects int) error {
	isChanged := this.server != server
	this.server = server

	if redirects > 8 {
		return errors.New("too many redirects")
	}
	redirects ++

	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return err
	}
	path := uri.Path

	// root
	if isChanged {
		this.root = server.Root
		if len(this.root) > 0 {
			this.root = this.format(this.root)
		}
	}

	// 字符集
	if len(server.Charset) > 0 {
		this.charset = this.format(server.Charset)
	}

	// Header
	if len(server.Headers) > 0 {
		this.headers = append(this.headers, server.FormatHeaders(func(source string) string {
			return this.format(source)
		}) ...)
	}

	if len(server.IgnoreHeaders) > 0 {
		this.ignoreHeaders = append(this.ignoreHeaders, server.IgnoreHeaders ...)
	}

	// location的相关配置
	for _, location := range server.Locations {
		if !location.On {
			continue
		}
		if locationMatches, ok := location.Match(path); ok {
			this.convertVarMapping(locationMatches)

			if len(location.Root) > 0 {
				this.root = this.format(location.Root)
			}
			if len(location.Charset) > 0 {
				this.charset = this.format(location.Charset)
			}
			if len(location.Index) > 0 {
				this.index = this.formatAll(location.Index)
			}
			if len(location.Headers) > 0 {
				this.headers = append(this.headers, location.FormatHeaders(func(source string) string {
					return this.format(source)
				}) ...)
			}

			if len(this.ignoreHeaders) > 0 {
				this.ignoreHeaders = append(this.ignoreHeaders, location.IgnoreHeaders ...)
			}

			this.location = location

			// rewrite相关配置
			if len(location.Rewrite) > 0 {
				for _, rule := range location.Rewrite {
					if !rule.On {
						continue
					}

					if replace, ok := rule.Match(path, func(source string) string {
						return this.format(source)
					}); ok {
						this.rewriteId = rule.Id

						if len(rule.Headers) > 0 {
							this.headers = append(this.headers, rule.FormatHeaders(func(source string) string {
								return this.format(source)
							}) ...)
						}

						if len(rule.IgnoreHeaders) > 0 {
							this.ignoreHeaders = append(this.ignoreHeaders, rule.IgnoreHeaders ...)
						}

						// @TODO 支持带host前缀的URL，比如：http://google.com/hello/world
						newURI, err := url.ParseRequestURI(replace)
						if err != nil {
							this.uri = replace
							return nil
						}
						if len(newURI.RawQuery) > 0 {
							this.uri = newURI.Path + "?" + newURI.RawQuery
							if len(uri.RawQuery) > 0 {
								this.uri += "&" + uri.RawQuery
							}
						} else {
							this.uri = newURI.Path
							if len(uri.RawQuery) > 0 {
								this.uri += "?" + uri.RawQuery
							}
						}

						switch rule.TargetType() {
						case teaconfigs.RewriteTargetURL:
							return this.configure(server, redirects)
						case teaconfigs.RewriteTargetProxy:
							proxyId := rule.TargetProxy()
							server, found := FindServer(proxyId)
							if !found {
								return errors.New("server with '" + proxyId + "' not found")
							}
							if !server.On {
								return errors.New("server with '" + proxyId + "' not available now")
							}
							return this.configure(server, redirects)
						}
						return nil
					}
				}
			}

			// fastcgi
			fastcgi := location.NextFastcgi()
			if fastcgi != nil {
				this.fastcgi = fastcgi

				if len(fastcgi.Headers) > 0 {
					this.headers = append(this.headers, fastcgi.Headers ...)
				}

				if len(fastcgi.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, fastcgi.IgnoreHeaders ...)
				}

				continue
			}

			// proxy
			if len(location.Proxy) > 0 {
				server, found := FindServer(location.Proxy)
				if !found {
					return errors.New("server with '" + location.Proxy + "' not found")
				}
				if !server.On {
					return errors.New("server with '" + location.Proxy + "' not available now")
				}
				return this.configure(server, redirects)
			}

			// backends
			if len(location.Backends) > 0 {
				backend := location.NextBackend()
				if backend == nil {
					return errors.New("no backends available")
				}
				this.backend = backend

				if len(backend.Headers) > 0 {
					this.headers = append(this.headers, backend.Headers ...)
				}

				if len(backend.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, backend.IgnoreHeaders ...)
				}

				continue
			}
		}
	}

	// server的相关配置
	if len(server.Rewrite) > 0 {
		for _, rule := range server.Rewrite {
			if !rule.On {
				continue
			}
			if replace, ok := rule.Match(path, func(source string) string {
				return this.format(source)
			}); ok {
				this.rewriteId = rule.Id

				if len(rule.Headers) > 0 {
					this.headers = append(this.headers, rule.Headers ...)
				}

				if len(rule.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, rule.IgnoreHeaders ...)
				}

				// @TODO 支持带host前缀的URL，比如：http://google.com/hello/world
				newURI, err := url.ParseRequestURI(replace)
				if err != nil {
					this.uri = replace
					return nil
				}
				if len(newURI.RawQuery) > 0 {
					this.uri = newURI.Path + "?" + newURI.RawQuery
					if len(uri.RawQuery) > 0 {
						this.uri += "&" + uri.RawQuery
					}
				} else {
					if len(uri.RawQuery) > 0 {
						this.uri = newURI.Path + "?" + uri.RawQuery
					}
				}

				switch rule.TargetType() {
				case teaconfigs.RewriteTargetURL:
					return this.configure(server, redirects)
				case teaconfigs.RewriteTargetProxy:
					proxyId := rule.TargetProxy()
					server, found := FindServer(proxyId)
					if !found {
						return errors.New("server with '" + proxyId + "' not found")
					}
					if !server.On {
						return errors.New("server with '" + proxyId + "' not available now")
					}
					return this.configure(server, redirects)
				}
				return nil
			}
		}
	}

	// fastcgi
	fastcgi := server.NextFastcgi()
	if fastcgi != nil {
		this.fastcgi = fastcgi

		if len(fastcgi.Headers) > 0 {
			this.headers = append(this.headers, fastcgi.Headers ...)
		}

		if len(fastcgi.IgnoreHeaders) > 0 {
			this.ignoreHeaders = append(this.ignoreHeaders, fastcgi.IgnoreHeaders ...)
		}

		return nil
	}

	// proxy
	if len(server.Proxy) > 0 {
		server, found := FindServer(server.Proxy)
		if !found {
			return errors.New("server with '" + server.Proxy + "' not found")
		}
		if !server.On {
			return errors.New("server with '" + server.Proxy + "' not available now")
		}
		return this.configure(server, redirects)
	}

	// 转发到后端
	backend := server.NextBackend()
	if backend == nil {
		if len(this.root) == 0 {
			return errors.New("no backends available")
		}
	}
	this.backend = backend

	if backend != nil {
		if len(backend.Headers) > 0 {
			this.headers = append(this.headers, backend.Headers ...)
		}

		if len(backend.IgnoreHeaders) > 0 {
			this.ignoreHeaders = append(this.ignoreHeaders, backend.IgnoreHeaders ...)
		}
	}

	return nil
}

func (this *Request) Call(writer http.ResponseWriter) error {
	defer this.log()

	if this.backend != nil {
		return this.callBackend(writer)
	}
	if this.proxy != nil {
		return this.callProxy(writer)
	}
	if this.fastcgi != nil {
		return this.callFastcgi(writer)
	}
	if len(this.root) > 0 {
		return this.callRoot(writer)
	}
	return errors.New("unable to handle the request")
}

// @TODO 支持eTag，cache等
func (this *Request) callRoot(writer http.ResponseWriter) error {
	if len(this.uri) == 0 {
		this.notFoundError(writer)
		return nil
	}

	requestPath := this.uri
	uri, err := url.ParseRequestURI(this.uri)
	query := ""
	if err == nil {
		requestPath = uri.Path
		query = uri.RawQuery
	}

	// 去掉其中的奇怪的路径
	requestPath = strings.Replace(requestPath, "..\\", "", -1)

	if requestPath == "/" {
		// 根目录
		indexFile := this.findIndexFile(this.root)
		if len(indexFile) > 0 {
			this.uri = requestPath + indexFile
			if len(query) > 0 {
				this.uri += "?" + query
			}
			err := this.configure(this.server, 0)
			if err != nil {
				logs.Error(err)
				this.serverError(writer)
				return nil
			}
			return this.Call(writer)
		} else {
			this.notFoundError(writer)
			return nil
		}
	}
	filename := strings.Replace(requestPath, "/", Tea.DS, -1)
	filePath := ""
	if filename[0:1] == Tea.DS {
		filePath = this.root + filename
	} else {
		filePath = this.root + Tea.DS + filename
	}

	this.filePath = filePath

	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			this.notFoundError(writer)
			return nil
		} else {
			this.serverError(writer)
			logs.Error(err)
			return nil
		}
	}
	if stat.IsDir() {
		indexFile := this.findIndexFile(filePath)
		if len(indexFile) > 0 {
			this.uri = requestPath + indexFile
			if len(query) > 0 {
				this.uri += "?" + query
			}
			err := this.configure(this.server, 0)
			if err != nil {
				logs.Error(err)
				this.serverError(writer)
				return nil
			}
			return this.Call(writer)
		} else {
			this.notFoundError(writer)
			return nil
		}
	}

	fp, err := os.OpenFile(filePath, os.O_RDONLY, 444)
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}
	defer fp.Close()

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// mime type
	if !hasIgnoreHeaders || !ignoreHeaders.Has("CONTENT-TYPE") {
		ext := filepath.Ext(requestPath)
		if len(ext) > 0 {
			mimeType := mime.TypeByExtension(ext)
			if len(mimeType) > 0 {
				if len(this.charset) > 0 {
					// 去掉里面的charset设置
					index := strings.Index(mimeType, "charset=")
					if index > 0 {
						writer.Header().Set("Content-Type", mimeType[:index+len("charset=")]+this.charset)
					} else {
						writer.Header().Set("Content-Type", mimeType+"; charset="+this.charset)
					}
				} else {
					writer.Header().Set("Content-Type", mimeType)
				}
			}
		}
	}

	// 自定义Header
	for _, header := range this.headers {
		if header.Match(http.StatusOK) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			writer.Header().Set(header.Name, header.Value)
		}
	}

	this.responseHeader = writer.Header()

	n, err := io.Copy(writer, fp)

	if err != nil {
		if this.debug {
			logs.Error(err)
		}
		return nil
	}

	this.responseStatus = http.StatusOK
	this.responseStatusMessage = "200 OK"
	this.responseBytesSent = n
	this.responseBodyBytesSent = n
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	return nil
}

func (this *Request) callBackend(writer http.ResponseWriter) error {
	if len(this.backend.Address) == 0 {
		this.serverError(writer)
		logs.Error(errors.New("backend address should not be empty"))
		return nil
	}

	this.raw.URL.Scheme = this.scheme
	this.raw.URL.Host = this.host

	// 设置代理相关的头部
	// 参考 https://tools.ietf.org/html/rfc7239
	this.raw.Header.Set("X-Real-IP", this.raw.RemoteAddr)
	this.raw.Header.Set("X-Forwarded-For", this.raw.RemoteAddr)
	this.raw.Header.Set("X-Forwarded-Host", this.host)
	this.raw.Header.Set("X-Forwarded-By", this.raw.RemoteAddr)
	this.raw.Header.Set("X-Forwarded-Proto", this.raw.Proto)
	//this.raw.Header.Set("Connection", "keep-alive")

	// @TODO 使用client池
	client := http.Client{
		Timeout: 30 * time.Second,

		// 处理跳转
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if via[0].URL.Host == req.URL.Host {
				http.Redirect(writer, this.raw, req.URL.RequestURI(), http.StatusTemporaryRedirect)
			} else {
				http.Redirect(writer, this.raw, req.URL.String(), http.StatusTemporaryRedirect)
			}
			return &RedirectError{}
		},
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 后端地址
			addr = this.backend.Address

			// 握手配置
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext(ctx, network, addr)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	this.raw.RequestURI = ""
	resp, err := client.Do(this.raw)
	if err != nil {
		urlError, ok := err.(*url.Error)
		if ok {
			if _, ok := urlError.Err.(*RedirectError); ok {
				return nil
			}
		}
		this.serverError(writer)
		logs.Error(err)
		return nil
	}
	defer resp.Body.Close()

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 设置Header
	hasCharset := len(this.charset) > 0
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(k)) {
			continue
		}
		for _, subV := range v {
			// 字符集
			if hasCharset && k == "Content-Type" {
				if _, found := textMimeMap[subV]; found {
					if !strings.Contains(subV, "charset=") {
						subV += "; charset=" + this.charset
					}
				}
			}

			writer.Header().Add(k, subV)
		}
	}

	// 自定义Header
	for _, header := range this.headers {
		if header.Match(resp.StatusCode) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			writer.Header().Set(header.Name, header.Value)
		}
	}

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	n, err := io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
		return nil
	}

	// 请求信息
	this.responseStatus = resp.StatusCode
	this.responseStatusMessage = resp.Status
	this.responseBytesSent = n
	this.responseBodyBytesSent = n
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	return nil
}

func (this *Request) callProxy(writer http.ResponseWriter) error {
	backend := this.proxy.NextBackend()
	this.backend = backend
	return this.callBackend(writer)
}

func (this *Request) callFastcgi(writer http.ResponseWriter) error {
	env := this.fastcgi.FilterParams(this.raw)
	if len(this.root) > 0 {
		if !env.Has("DOCUMENT_ROOT") {
			env["DOCUMENT_ROOT"] = this.root
		}
	}
	if !env.Has("REMOTE_ADDR") {
		env["REMOTE_ADDR"] = this.raw.RemoteAddr
	}
	if !env.Has("QUERY_STRING") {
		u, err := url.ParseRequestURI(this.uri)
		if err == nil {
			env["QUERY_STRING"] = u.RawQuery
		} else {
			env["QUERY_STRING"] = this.raw.URL.RawQuery
		}
	}
	if !env.Has("SERVER_NAME") {
		env["SERVER_NAME"] = this.host
	}
	if !env.Has("REQUEST_URI") {
		env["REQUEST_URI"] = this.uri
	}
	if !env.Has("HOST") {
		env["HOST"] = this.host
	}

	if len(this.serverAddr) > 0 {
		if !env.Has("SERVER_ADDR") {
			env["SERVER_ADDR"] = this.serverAddr
		}
		if !env.Has("SERVER_PORT") {
			portIndex := strings.LastIndex(this.serverAddr, ":")
			if portIndex >= 0 {
				env["SERVER_PORT"] = this.serverAddr[portIndex+1:]
			}
		}
	}

	// 连接池配置
	poolSize := this.fastcgi.PoolSize
	if poolSize <= 0 {
		poolSize = 16
	}

	client, err := gofcgi.SharedPool(this.fastcgi.Network(), this.fastcgi.Address(), uint(poolSize)).Client()
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}

	// 请求相关
	if !env.Has("REQUEST_METHOD") {
		env["REQUEST_METHOD"] = this.method
	}
	if !env.Has("CONTENT_LENGTH") {
		env["CONTENT_LENGTH"] = fmt.Sprintf("%d", this.raw.ContentLength)
	}
	if !env.Has("CONTENT_TYPE") {
		env["CONTENT_TYPE"] = this.raw.Header.Get("Content-Type")
	}

	// 处理SCRIPT_FILENAME
	scriptFilename := env.GetString("SCRIPT_FILENAME")
	if len(scriptFilename) > 0 && (strings.Index(scriptFilename, "/") < 0 && strings.Index(scriptFilename, "\\") < 0) {
		env["SCRIPT_FILENAME"] = env.GetString("DOCUMENT_ROOT") + Tea.DS + scriptFilename
	}

	params := map[string]string{}
	for key, value := range env {
		params[key] = types.String(value)
	}

	for k, v := range this.raw.Header {
		if k == "Connection" {
			continue
		}
		for _, subV := range v {
			params["HTTP_"+strings.ToUpper(strings.Replace(k, "-", "_", -1))] = subV
		}
	}

	host, found := params["HTTP_HOST"]
	if !found || len(host) == 0 {
		params["HTTP_HOST"] = this.host
	}

	fcgiReq := gofcgi.NewRequest()
	fcgiReq.SetTimeout(this.fastcgi.Timeout())
	fcgiReq.SetParams(params)
	fcgiReq.SetBody(this.raw.Body, uint32(this.requestLength()))

	resp, err := client.Call(fcgiReq)
	if err != nil {
		this.serverError(writer)
		//if this.debug {
		logs.Error(err)
		//}
		return nil
	}

	defer resp.Body.Close()

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 设置Header
	var hasCharset = len(this.charset) > 0
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(k)) {
			continue
		}

		for _, subV := range v {
			// 字符集
			if hasCharset && k == "Content-Type" {
				if _, found := textMimeMap[subV]; found {
					if !strings.Contains(subV, "charset=") {
						subV += "; charset=" + this.charset
					}
				}
			}
			writer.Header().Add(k, subV)
		}
	}

	// 自定义Header
	for _, header := range this.headers {
		if header.Match(resp.StatusCode) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			writer.Header().Set(header.Name, header.Value)
		}
	}

	// 插件过滤
	if teaplugins.HasResponseFilters {
		resp.Header = writer.Header()
		resp = teaplugins.FilterResponse(resp)

		// reset headers
		oldHeaders := writer.Header()
		for key := range oldHeaders {
			oldHeaders.Del(key)
		}

		for key, value := range resp.Header {
			for _, v := range value {
				oldHeaders.Add(key, v)
			}
		}
	}

	this.responseHeader = writer.Header()

	// 设置响应码
	writer.WriteHeader(resp.StatusCode)

	// 输出内容
	n, err := io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
		return nil
	}

	// 请求信息
	this.responseStatus = resp.StatusCode
	this.responseStatusMessage = resp.Status
	this.responseBytesSent = n
	this.responseBodyBytesSent = n
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	return nil
}

func (this *Request) notFoundError(writer http.ResponseWriter) {
	msg := "404 page not found: '" + this.requestURI() + "'"

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(msg))

	this.responseStatus = http.StatusNotFound
	this.responseStatusMessage = msg
	this.responseBodyBytesSent = int64(len(msg))
	this.responseBytesSent = this.responseBodyBytesSent
}

func (this *Request) serverError(writer http.ResponseWriter) {
	msg := "500 Internal Server Error"

	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(msg))

	this.responseStatus = http.StatusInternalServerError
	this.responseStatusMessage = msg
	this.responseBodyBytesSent = int64(len(msg))
	this.responseBytesSent = this.responseBodyBytesSent
}

func (this *Request) requestRemoteAddr() string {
	// Real-IP
	realIP := this.raw.Header.Get("X-Real-IP")
	if len(realIP) > 0 {
		index := strings.LastIndex(realIP, ":")
		if index < 0 {
			return realIP
		} else {
			return realIP[:index]
		}
	}

	// X-Forwarded-For
	forwardedFor := this.raw.Header.Get("X-Forwarded-For")
	if len(forwardedFor) > 0 {
		index := strings.LastIndex(forwardedFor, ":")
		if index < 0 {
			return forwardedFor
		} else {
			return forwardedFor[:index]
		}
	}

	// Remote-Addr
	remoteAddr := this.raw.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index < 0 {
		return remoteAddr
	} else {
		return remoteAddr[:index]
	}
}

func (this *Request) requestRemotePort() int {
	remoteAddr := this.raw.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index < 0 {
		return 0
	} else {
		return types.Int(remoteAddr[index+1:])
	}
}

func (this *Request) requestRemoteUser() string {
	username, _, ok := this.raw.BasicAuth()
	if !ok {
		return ""
	}
	return username
}

func (this *Request) requestURI() string {
	return this.rawURI
}

func (this *Request) requestPath() string {
	uri, err := url.ParseRequestURI(this.requestURI())
	if err != nil {
		return ""
	}
	return uri.Path
}

func (this *Request) requestLength() int64 {
	return this.raw.ContentLength
}

func (this *Request) requestMethod() string {
	return this.method
}

func (this *Request) requestFilename() string {
	return this.filePath
}

func (this *Request) requestProto() string {
	return this.raw.Proto
}

func (this *Request) requestReferer() string {
	return this.raw.Referer()
}

func (this *Request) requestUserAgent() string {
	return this.raw.UserAgent()
}

func (this *Request) requestContentType() string {
	return this.raw.Header.Get("Content-Type")
}

func (this *Request) requestString() string {
	return this.method + " " + this.requestURI() + " " + this.requestProto()
}

func (this *Request) requestCookiesString() string {
	var cookies = []string{}
	for _, cookie := range this.raw.Cookies() {
		cookies = append(cookies, url.QueryEscape(cookie.Name)+"="+url.QueryEscape(cookie.Value))
	}
	return strings.Join(cookies, "&")
}

func (this *Request) requestCookie(name string) string {
	cookie, err := this.raw.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Name
}

func (this *Request) requestQueryString() string {
	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return ""
	}
	return uri.RawQuery
}

func (this *Request) requestQueryParam(name string) string {
	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return ""
	}

	v, found := uri.Query()[name]
	if !found {
		return ""
	}
	return strings.Join(v, "&")
}

func (this *Request) requestServerPort() int {
	index := strings.LastIndex(this.serverAddr, ":")
	if index < 0 {
		return 0
	}
	return types.Int(this.serverAddr[index+1:])
}

func (this *Request) requestHeadersString() string {
	var headers = []string{}
	for k, v := range this.raw.Header {
		for _, subV := range v {
			headers = append(headers, k+": "+subV)
		}
	}
	return strings.Join(headers, ";")
}

func (this *Request) requestHeader(key string) string {
	v, found := this.raw.Header[key]
	if !found {
		return ""
	}
	return strings.Join(v, ";")
}

// 利用请求参数格式化字符串
func (this *Request) format(source string) string {
	if len(source) == 0 {
		return ""
	}

	var varName = ""
	var hasVarMapping = len(this.varMapping) > 0
	return requestVarReg.ReplaceAllStringFunc(source, func(s string) string {
		varName = s[2 : len(s)-1]

		// 自定义变量
		if hasVarMapping {
			value, found := this.varMapping[varName]
			if found {
				return value
			}
		}

		// 请求变量
		switch varName {
		case "teaVersion":
			return teaconst.TeaVersion
		case "remoteAddr":
			return this.requestRemoteAddr()
		case "remotePort":
			return fmt.Sprintf("%d", this.requestRemotePort())
		case "remoteUser":
			return this.requestRemoteUser()
		case "requestURI", "requestUri":
			return this.requestURI()
		case "requestPath":
			return this.requestPath()
		case "requestLength":
			return fmt.Sprintf("%d", this.requestLength())
		case "requestTime":
			return fmt.Sprintf("%.6f", this.requestTime)
		case "requestMethod":
			return this.requestMethod()
		case "requestFilename":
			return this.requestFilename()
		case "scheme":
			return this.rawScheme
		case "serverProtocol", "proto":
			return this.requestProto()
		case "bytesSent":
			return fmt.Sprintf("%d", this.responseBytesSent)
		case "bodyBytesSent":
			return fmt.Sprintf("%d", this.responseBodyBytesSent)
		case "status":
			return fmt.Sprintf("%d", this.responseStatus)
		case "statusMessage":
			return this.responseStatusMessage
		case "timeISO8601":
			return this.requestTimeISO8601
		case "timeLocal":
			return this.requestTimeLocal
		case "msec":
			return fmt.Sprintf("%.6f", this.requestMsec)
		case "timestamp":
			return fmt.Sprintf("%d", this.requestTimestamp)
		case "host":
			return this.host
		case "referer":
			return this.requestReferer()
		case "userAgent":
			return this.requestUserAgent()
		case "contentType":
			return this.requestContentType()
		case "request":
			return this.requestString()
		case "cookies":
			return this.requestCookiesString()
		case "args", "queryString":
			return this.requestQueryString()
		case "headers":
			return this.requestHeadersString()
		case "serverName":
			return this.serverName
		case "serverPort":
			return fmt.Sprintf("%d", this.requestServerPort())
		}

		dotIndex := strings.Index(varName, ".")
		if dotIndex < 0 {
			return s
		}
		prefix := varName[:dotIndex]
		suffix := varName[dotIndex+1:]

		// cookie.
		if prefix == "cookie" {
			return this.requestCookie(suffix)
		}

		// arg.
		if prefix == "arg" {
			return this.requestQueryParam(suffix)
		}

		// header.
		if prefix == "header" || prefix == "http" {
			return this.requestHeader(suffix)
		}

		return s
	})
}

// 格式化一组字符串
func (this *Request) formatAll(sources []string) []string {
	result := []string{}
	for _, s := range sources {
		result = append(result, this.format(s))
	}
	return result
}

// 记录日志
func (this *Request) log() {
	if !this.shouldLog {
		return
	}

	cookies := map[string]string{}
	for _, cookie := range this.raw.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	accessLog := &tealogs.AccessLog{
		TeaVersion:      teaconst.TeaVersion,
		RemoteAddr:      this.requestRemoteAddr(),
		RemotePort:      this.requestRemotePort(),
		RemoteUser:      this.requestRemoteUser(),
		RequestURI:      this.requestURI(),
		RequestPath:     this.requestPath(),
		RequestLength:   this.requestLength(),
		RequestTime:     this.requestTime,
		RequestMethod:   this.requestMethod(),
		RequestFilename: this.requestFilename(),
		Scheme:          this.scheme,
		Proto:           this.requestProto(),
		BytesSent:       this.responseBytesSent,
		BodyBytesSent:   this.responseBodyBytesSent,
		Status:          this.responseStatus,
		StatusMessage:   this.responseStatusMessage,
		TimeISO8601:     this.requestTimeISO8601,
		TimeLocal:       this.requestTimeLocal,
		Msec:            this.requestMsec,
		Timestamp:       this.requestTimestamp,
		Host:            this.host,
		Referer:         this.requestReferer(),
		UserAgent:       this.requestUserAgent(),
		Request:         this.requestString(),
		ContentType:     this.requestContentType(),
		Cookie:          cookies,
		Args:            this.requestQueryString(),
		QueryString:     this.requestQueryString(),
		Header:          this.raw.Header,
		ServerName:      this.serverName,
		ServerPort:      this.requestServerPort(),
		ServerProtocol:  this.requestProto(),
	}

	if this.server != nil {
		accessLog.ServerId = this.server.Id
	}

	if this.backend != nil {
		accessLog.BackendAddress = this.backend.Address
		accessLog.BackendId = this.backend.Id
	}

	if this.fastcgi != nil {
		accessLog.FastcgiAddress = this.fastcgi.Pass
		accessLog.FastcgiId = this.fastcgi.Id
	}

	accessLog.RewriteId = this.rewriteId

	if this.location != nil {
		accessLog.LocationId = this.location.Id
	}

	if this.responseHeader != nil {
		accessLog.SentHeader = this.responseHeader
	}

	tealogs.SharedLogger().Push(accessLog)
}

func (this *Request) findIndexFile(dir string) string {
	if len(this.index) == 0 {
		return ""
	}
	for _, index := range this.index {
		if len(index) == 0 {
			continue
		}

		// 模糊查找
		if strings.Contains(index, "*") {
			files, err := filepath.Glob(dir + Tea.DS + index)
			if err != nil {
				logs.Error(err)
				continue
			}
			if len(files) > 0 {
				return filepath.Base(files[0])
			}
			continue
		}

		// 精确查找
		filePath := dir + Tea.DS + index
		stat, err := os.Stat(filePath)
		if err != nil || !stat.Mode().IsRegular() {
			continue
		}
		return index
	}
	return ""
}

func (this *Request) convertIgnoreHeaders() maps.Map {
	m := maps.Map{}
	for _, h := range this.ignoreHeaders {
		m[strings.ToUpper(h)] = true
	}
	return m
}

func (this *Request) convertVarMapping(s []string) {
	for k, v := range s {
		this.varMapping[fmt.Sprintf("%d", k)] = v
	}
}
