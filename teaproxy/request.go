package teaproxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/TeaWeb/code/teautils"
	"github.com/gorilla/websocket"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/gofcgi"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	apiconfig "github.com/TeaWeb/code/teaconfigs/api"
)

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
	headers       []*shared.HeaderConfig // 自定义Header
	ignoreHeaders []string               // 忽略的Header
	varMapping    map[string]string      // 自定义变量

	root         string   // 资源根目录
	index        []string // 目录下默认访问的文件
	backend      *teaconfigs.BackendConfig
	fastcgi      *teaconfigs.FastcgiConfig
	proxy        *teaconfigs.ServerConfig
	location     *teaconfigs.LocationConfig
	accessPolicy *shared.AccessPolicy

	cachePolicy  *shared.CachePolicy
	cacheEnabled bool

	api    *apiconfig.API // API
	mockOn bool           // 是否开启了API Mock

	rewriteId           string // 匹配的rewrite id
	rewriteReplace      string // 经过rewrite之后的URL
	rewriteRedirectMode string // 跳转方式
	rewriteIsExternal   bool   // 是否为外部URL

	websocket *teaconfigs.WebsocketConfig

	// 执行请求
	filePath string

	responseWriter   *ResponseWriter
	responseCallback func(http.ResponseWriter)

	requestFromTime    time.Time // 请求开始时间
	requestTime        float64   // 请求耗时
	requestTimeISO8601 string
	requestTimeLocal   string
	requestMsec        float64
	requestTimestamp   int64
	requestMaxSize     int64

	isWatching        bool     // 是否在监控
	requestData       []byte   // 导出的request，在监控请求的时候有用
	responseAPIStatus string   // API状态码
	errors            []string // 错误信息

	enableAccessLog bool
	gzipLevel       uint8
	gzipMinLength   int64
	debug           bool
}

// 获取新的请求
func NewRequest(rawRequest *http.Request) *Request {
	now := time.Now()
	req := &Request{
		varMapping:         map[string]string{},
		raw:                rawRequest,
		rawURI:             rawRequest.URL.RequestURI(),
		requestFromTime:    now,
		requestTimestamp:   now.Unix(),
		requestTimeISO8601: now.Format("2006-01-02T15:04:05.000Z07:00"),
		requestTimeLocal:   now.Format("2/Jan/2006:15:04:05 -0700"),
		requestMsec:        float64(now.Unix()) + float64(now.Nanosecond())/1000000000,
		enableAccessLog:    true,
	}

	return req
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
			this.root = this.Format(this.root)
		}
	}

	// 字符集
	if len(server.Charset) > 0 {
		this.charset = this.Format(server.Charset)
	}

	// Header
	if len(server.Headers) > 0 {
		// 延迟执行，让Header有机会加入Backend, Fastcgi等信息
		defer func() {
			this.headers = append(this.headers, server.FormatHeaders(this.Format) ...)
		}()
	}

	if len(server.IgnoreHeaders) > 0 {
		this.ignoreHeaders = append(this.ignoreHeaders, server.IgnoreHeaders ...)
	}

	// cache
	if server.CacheOn {
		cachePolicy := server.CachePolicyObject()
		if cachePolicy != nil && cachePolicy.On {
			this.cachePolicy = cachePolicy
		}
	}

	// other
	if server.MaxBodyBytes() > 0 {
		this.requestMaxSize = server.MaxBodyBytes()
	}
	if server.DisableAccessLog {
		this.enableAccessLog = false
	}
	this.gzipLevel = server.GzipLevel
	this.gzipMinLength = server.GzipMinBytes()

	// API配置，目前只有Plus版本支持
	if teaconst.PlusEnabled && server.API != nil && server.API.On {
		// 查找API
		api, params := server.API.FindActiveAPI(uri.Path, this.method)
		if api != nil {
			this.api = api

			// cache
			if api.CacheOn {
				cachePolicy := api.CachePolicyObject()
				if cachePolicy != nil && cachePolicy.On {
					this.cachePolicy = cachePolicy
				}
			}

			// address
			if len(api.Address) > 0 {
				address := api.Address

				query := uri.Query()

				// 匹配的参数
				if params != nil {
					for key, value := range params {
						query[key] = []string{value}
					}
					uri.RawQuery = query.Encode()
				}

				// 支持变量
				address = teaconfigs.RegexpNamedVariable.ReplaceAllStringFunc(address, func(s string) string {
					match := s[2 : len(s)-1]
					switch match {
					case "path":
						return api.Path
					}

					if strings.HasPrefix(match, "arg.") {
						value, found := query[match[4:]]
						if found {
							return strings.Join(value, ",")
						} else {
							return ""
						}
					}

					return s
				})

				this.uri = address

				newURI, err := url.ParseRequestURI(address)
				if err == nil {
					path = newURI.Path

					// query
					if len(uri.RawQuery) > 0 {
						if len(newURI.RawQuery) == 0 {
							this.uri = address + "?" + uri.RawQuery
						} else {
							this.uri = address + "&" + uri.RawQuery
						}
					}
				} else {
					path = address
				}
			}

			// headers
			if len(api.Headers) > 0 {
				this.headers = append(this.headers, api.FormatHeaders(func(source string) string {
					return this.Format(source)
				}) ...)
			}

			// dump
			if api.IsWatching() {
				this.isWatching = true

				// 判断如果Content-Length过长，则截断
				reqData, err := httputil.DumpRequest(this.raw, true)
				if err == nil {
					if len(reqData) > 10240 {
						reqData = reqData[:10240]
					}
					this.requestData = reqData
				}
			}

			// 是否有Mock
			if server.API.MockOn && api.MockOn && len(api.MockFiles) > 0 {
				this.mockOn = true
				return nil
			}
		}
	}

	// location的相关配置
	var locationConfigured = false
	for _, location := range server.Locations {
		if !location.On {
			continue
		}
		if locationMatches, ok := location.Match(path); ok {
			this.addVarMapping(locationMatches)

			if len(location.Root) > 0 {
				this.root = this.Format(location.Root)
				locationConfigured = true
			}
			if len(location.Charset) > 0 {
				this.charset = this.Format(location.Charset)
			}
			if len(location.Index) > 0 {
				this.index = this.formatAll(location.Index)
			}
			if location.MaxBodyBytes() > 0 {
				this.requestMaxSize = location.MaxBodyBytes()
			}
			if location.DisableAccessLog {
				this.enableAccessLog = false
			}
			if location.GzipLevel >= 0 {
				this.gzipLevel = uint8(location.GzipLevel)
			}
			if location.GzipMinBytes() > 0 {
				this.gzipMinLength = location.GzipMinBytes()
			}

			if location.CacheOn {
				cachePolicy := location.CachePolicyObject()
				if cachePolicy != nil && cachePolicy.On {
					this.cachePolicy = cachePolicy
				}
			}

			if len(location.Headers) > 0 {
				this.headers = append(this.headers, location.FormatHeaders(this.Format) ...)
			}

			if len(location.IgnoreHeaders) > 0 {
				this.ignoreHeaders = append(this.ignoreHeaders, location.IgnoreHeaders ...)
			}

			if location.AccessPolicy != nil {
				this.accessPolicy = location.AccessPolicy
			}

			this.location = location

			// rewrite相关配置
			if len(location.Rewrite) > 0 {
				for _, rule := range location.Rewrite {
					if !rule.On {
						continue
					}

					if replace, varMapping, ok := rule.Match(path, this.Format); ok {
						this.addVarMapping(varMapping)
						this.rewriteId = rule.Id

						if len(rule.Headers) > 0 {
							this.headers = append(this.headers, rule.FormatHeaders(func(source string) string {
								return this.Format(source)
							}) ...)
						}

						if len(rule.IgnoreHeaders) > 0 {
							this.ignoreHeaders = append(this.ignoreHeaders, rule.IgnoreHeaders ...)
						}

						// 外部URL
						if rule.IsExternalURL(replace) {
							this.rewriteReplace = replace
							this.rewriteIsExternal = true
							this.rewriteRedirectMode = rule.RedirectMode()
							return nil
						}

						// 内部URL
						if rule.RedirectMode() == teaconfigs.RewriteFlagRedirect {
							this.rewriteReplace = replace
							this.rewriteIsExternal = false
							this.rewriteRedirectMode = teaconfigs.RewriteFlagRedirect
							return nil
						}

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
				this.backend = nil // 防止冲突
				locationConfigured = true

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
				options := maps.Map{
					"request":   this.raw,
					"formatter": this.Format,
				}
				backend := location.NextBackend(options)
				if backend == nil {
					return errors.New("no backends available")
				}
				this.backend = backend
				locationConfigured = true

				if len(backend.Headers) > 0 {
					this.headers = append(this.headers, backend.Headers ...)
				}

				if len(backend.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, backend.IgnoreHeaders ...)
				}

				continue
			}

			// websocket
			if location.Websocket != nil && location.Websocket.On {
				options := maps.Map{
					"request":   this.raw,
					"formatter": this.Format,
				}
				this.backend = location.Websocket.NextBackend(options)
				this.websocket = location.Websocket
				return nil
			}
		}
	}

	// 如果经过location找到了相关配置，就终止
	if locationConfigured {
		return nil
	}

	// server的相关配置
	if len(server.Rewrite) > 0 {
		for _, rule := range server.Rewrite {
			if !rule.On {
				continue
			}
			if replace, varMapping, ok := rule.Match(path, func(source string) string {
				return this.Format(source)
			}); ok {
				this.addVarMapping(varMapping)
				this.rewriteId = rule.Id

				if len(rule.Headers) > 0 {
					this.headers = append(this.headers, rule.Headers ...)
				}

				if len(rule.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, rule.IgnoreHeaders ...)
				}

				// 外部URL
				if rule.IsExternalURL(replace) {
					this.rewriteReplace = replace
					this.rewriteIsExternal = true
					this.rewriteRedirectMode = rule.RedirectMode()
					return nil
				}

				// 内部URL
				if rule.RedirectMode() == teaconfigs.RewriteFlagRedirect {
					this.rewriteReplace = replace
					this.rewriteIsExternal = false
					this.rewriteRedirectMode = teaconfigs.RewriteFlagRedirect
					return nil
				}

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
		this.backend = nil // 防止冲突

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
	options := maps.Map{
		"request":   this.raw,
		"formatter": this.Format,
	}
	backend := server.NextBackend(options)
	if backend == nil {
		if len(this.root) == 0 {
			return errors.New("no backends available")
		}
	}
	responseCallback := options.Get("responseCallback")
	if responseCallback != nil {
		f, ok := responseCallback.(func(http.ResponseWriter))
		if ok {
			this.responseCallback = f
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

func (this *Request) call(writer *ResponseWriter) error {
	defer func() {
		// log
		this.log()

		// call hook
		CallRequestAfterHook(this, writer)
	}()

	if this.requestMaxSize > 0 {
		this.raw.Body = http.MaxBytesReader(writer, this.raw.Body, this.requestMaxSize)
	}

	if this.gzipLevel > 0 && this.allowGzip() {
		writer.Gzip(this.gzipLevel, this.gzipMinLength)
		defer writer.Close()
	}

	// UV
	/**uid, err := this.raw.Cookie("TeaUID")
	if err != nil || uid == nil {
		http.SetCookie(writer, &http.Cookie{
			Name:    "TeaUID",
			Value:   stringutil.Rand(32),
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
	}**/

	this.responseWriter = writer

	// hook
	b := CallRequestBeforeHook(this, writer)
	if !b {
		return nil
	}

	// 是否有mock
	if this.mockOn {
		return this.callMock(writer)
	}

	// watch
	if this.isWatching {
		// 判断如果Content-Length过长，则截断
		reqData, err := httputil.DumpRequest(this.raw, true)
		if err == nil {
			if len(reqData) > 10240 {
				reqData = reqData[:10240]
			}
			this.requestData = reqData
		}

		writer.SetBodyCopying(true)
	}

	// API相关
	if this.api != nil {
		// limit
		if this.api.Limit != nil {
			this.api.Limit.Begin()
			defer this.api.Limit.Done()
		}

		// 检查consumer
		goNext := this.consumeAPI(writer)
		if !goNext {
			return nil
		}
	}

	// access policy
	if this.accessPolicy != nil {
		if !this.accessPolicy.AllowAccess(this.requestRemoteAddr()) {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Forbidden Request"))
			return nil
		}

		reason, allowed := this.accessPolicy.AllowTraffic()
		if !allowed {
			writer.WriteHeader(http.StatusTooManyRequests)
			writer.Write([]byte("[" + reason + "]Request Quota Exceeded"))
			return nil
		}
	}

	if this.websocket != nil {
		return this.callWebsocket(writer)
	}
	if this.backend != nil {
		return this.callBackend(writer)
	}
	if this.proxy != nil {
		return this.callProxy(writer)
	}
	if this.fastcgi != nil {
		return this.callFastcgi(writer)
	}
	if len(this.rewriteId) > 0 && (this.rewriteIsExternal || this.rewriteRedirectMode == teaconfigs.RewriteFlagRedirect) {
		return this.callRewrite(writer)
	}
	if len(this.root) > 0 {
		return this.callRoot(writer)
	}
	return errors.New("unable to handle the request")
}

// 调用本地静态资源
func (this *Request) callRoot(writer *ResponseWriter) error {
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
				this.addError(err)
				this.serverError(writer)
				return nil
			}
			return this.call(writer)
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
			this.addError(err)
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
				this.addError(err)
				return nil
			}
			return this.call(writer)
		} else {
			this.notFoundError(writer)
			return nil
		}
	}

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 响应header
	respHeader := writer.Header()

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
						respHeader.Set("Content-Type", mimeType[:index+len("charset=")]+this.charset)
					} else {
						respHeader.Set("Content-Type", mimeType+"; charset="+this.charset)
					}
				} else {
					respHeader.Set("Content-Type", mimeType)
				}
			}
		}
	}

	// length
	respHeader.Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	// 自定义Header
	for _, header := range this.headers {
		if header.Match(http.StatusOK) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			respHeader.Set(header.Name, header.Value)
		}
	}

	// 支持 Last-Modified
	modifiedTime := stat.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	if len(respHeader.Get("Last-Modified")) == 0 {
		respHeader.Set("Last-Modified", modifiedTime)
	}

	// 支持 ETag
	eTag := "\"et" + stringutil.Md5(fmt.Sprintf("%d,%d", stat.ModTime().UnixNano(), stat.Size())) + "\""
	if len(respHeader.Get("ETag")) == 0 {
		respHeader.Set("ETag", eTag)
	}

	// proxy callback
	if this.responseCallback != nil {
		this.responseCallback(writer)
	}

	// 支持 If-None-Match
	if this.requestHeader("If-None-Match") == eTag {
		writer.WriteHeader(http.StatusNotModified)

		return nil
	}

	// 支持 If-Modified-Since
	if this.requestHeader("If-Modified-Since") == modifiedTime {
		writer.WriteHeader(http.StatusNotModified)

		return nil
	}

	fp, err := os.OpenFile(filePath, os.O_RDONLY, 444)
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		this.addError(err)
		return nil
	}
	defer fp.Close()

	writer.Prepare(stat.Size())
	_, err = io.Copy(writer, fp)

	if err != nil {
		if this.debug {
			logs.Error(err)
		}
		return nil
	}

	return nil
}

// 调用Websocket
func (this *Request) callWebsocket(writer *ResponseWriter) error {
	if this.backend == nil {
		err := errors.New(this.requestPath() + ": no available backends for websocket")
		logs.Error(err)
		this.addError(err)
		this.serverError(writer)
		return err
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: this.websocket.HandshakeTimeoutDuration(),
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if len(origin) == 0 {
				return false
			}
			return this.websocket.MatchOrigin(origin)
		},
	}

	// 接收客户端连接
	client, err := upgrader.Upgrade(this.responseWriter.Raw(), this.raw, nil)
	if err != nil {
		logs.Error(errors.New("upgrade: " + err.Error()))
		this.addError(errors.New("upgrade: " + err.Error()))
		return err
	}
	defer client.Close()

	if this.websocket.ForwardMode == teaconfigs.WebsocketForwardModeWebsocket {
		// 判断最大连接数
		if this.backend.CurrentConns >= this.backend.MaxConns {
			this.serverError(writer)
			logs.Error(errors.New("too many connections"))
			this.addError(errors.New("too many connections"))
			return nil
		}

		// 增加连接数
		this.backend.IncreaseConn()
		defer this.backend.DecreaseConn()

		// 连接后端服务器
		wsURL := url.URL{Scheme: "ws", Host: this.backend.Address, Path: this.raw.RequestURI}
		dialer := websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: this.backend.FailTimeoutDuration(),
		}
		server, _, err := dialer.Dial(wsURL.String(), nil)
		if err != nil {
			logs.Error(err)
			this.addError(err)
			currentFails := this.backend.IncreaseFails()
			if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
				this.backend.IsDown = true
				this.backend.DownTime = time.Now()
				this.websocket.SetupScheduling(false)
			}
			return err
		}
		defer server.Close()

		// 设置关闭连接的处理函数
		clientIsClosed := false
		serverIsClosed := false
		client.SetCloseHandler(func(code int, text string) error {
			if serverIsClosed {
				return nil
			}
			serverIsClosed = true
			return server.Close()
		})

		// 从客户端接收数据
		go func() {
			for {
				messageType, message, err := client.ReadMessage()
				if err != nil {
					closeErr, ok := err.(*websocket.CloseError)
					if !ok || closeErr.Code != websocket.CloseGoingAway {
						logs.Error(err)
						this.addError(err)
					}
					clientIsClosed = true
					break
				}
				server.WriteMessage(messageType, message)
			}
		}()

		// 从后端服务器读取数据
		for {
			messageType, message, err := server.ReadMessage()
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
				if !ok || closeErr.Code != websocket.CloseGoingAway {
					logs.Error(err)
					this.addError(err)
				}
				serverIsClosed = true
				server.Close()
				if !clientIsClosed {
					client.Close()
				}
				break
			}
			client.WriteMessage(messageType, message)
		}
	} else if this.websocket.ForwardMode == teaconfigs.WebsocketForwardModeHttp {
		messageQueue := make(chan []byte, 1024)
		quit := make(chan bool)
		go func() {
		FOR:
			for {
				select {
				case message := <-messageQueue:
					{
						this.raw.Method = http.MethodPut
						responseWriter := NewResponseWriter(nil)
						responseWriter.SetBodyCopying(true)
						this.raw.Body = ioutil.NopCloser(bytes.NewReader(message))
						this.raw.Header.Del("Upgrade")
						err := this.callBackend(responseWriter)
						if err != nil {
							continue FOR
						}
						if responseWriter.StatusCode() != http.StatusOK {
							logs.Error(errors.New(this.requestURI() + ": invalid response from backend: " + fmt.Sprintf("%d", responseWriter.StatusCode()) + " " + http.StatusText(responseWriter.StatusCode())))
							this.addError(errors.New(this.requestURI() + ": invalid response from backend: " + fmt.Sprintf("%d", responseWriter.StatusCode()) + " " + http.StatusText(responseWriter.StatusCode())))
							continue FOR
						}
						client.WriteMessage(websocket.TextMessage, responseWriter.Body())
					}
				case <-quit:
					break FOR
				}
			}
		}()
		for {
			messageType, message, err := client.ReadMessage()
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
				if !ok || closeErr.Code != websocket.CloseGoingAway {
					logs.Error(err)
					this.addError(err)
				}
				quit <- true
				break
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				messageQueue <- message
			}
		}
	}

	return nil
}

// 调用后端服务器
func (this *Request) callBackend(writer *ResponseWriter) error {
	this.backend.IncreaseConn()
	defer this.backend.DecreaseConn()

	if len(this.backend.Address) == 0 {
		this.serverError(writer)
		logs.Error(errors.New("backend address should not be empty"))
		this.addError(errors.New("backend address should not be empty"))
		return nil
	}

	this.raw.URL.Host = this.host

	if len(this.backend.Scheme) > 0 && this.backend.Scheme != "http" {
		this.raw.URL.Scheme = this.backend.Scheme
	} else {
		this.raw.URL.Scheme = this.scheme
	}

	// new uri
	u, err := url.ParseRequestURI(this.uri)
	if err == nil {
		this.raw.URL.Path = u.Path
		this.raw.URL.RawQuery = u.RawQuery
	}

	// 设置代理相关的头部
	// 参考 https://tools.ietf.org/html/rfc7239
	if len(this.raw.RemoteAddr) > 0 {
		index := strings.Index(this.raw.RemoteAddr, ":")
		ip := ""
		if index > -1 {
			ip = this.raw.RemoteAddr[:index]
		} else {
			ip = this.raw.RemoteAddr
		}
		this.raw.Header.Set("X-Real-IP", ip)
		this.raw.Header.Set("X-Forwarded-For", ip)
		this.raw.Header.Set("X-Forwarded-By", ip)
	}
	this.raw.Header.Set("X-Forwarded-Host", this.host)
	this.raw.Header.Set("X-Forwarded-Proto", this.raw.Proto)

	client := SharedClientPool.client(this.backend.Id, this.backend.Address, this.backend.FailTimeoutDuration(), this.backend.ReadTimeoutDuration(), this.backend.MaxConns)

	this.raw.RequestURI = ""

	resp, err := client.Do(this.raw)
	if err != nil {
		urlError, ok := err.(*url.Error)
		if ok {
			if _, ok := urlError.Err.(*RedirectError); ok {
				http.Redirect(writer, this.raw, resp.Header.Get("Location"), resp.StatusCode)
				return nil
			}
		}

		// 如果超过最大失败次数，则下线
		currentFails := this.backend.IncreaseFails()
		if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
			this.backend.IsDown = true
			this.backend.DownTime = time.Now()
			if this.websocket != nil {
				this.websocket.SetupScheduling(false)
			} else {
				this.server.SetupScheduling(false)
			}
		}

		this.serverError(writer)
		logs.Error(err)
		this.addError(err)
		return nil
	}
	defer resp.Body.Close()

	// 清除错误次数
	if resp.StatusCode >= 200 {
		if !this.backend.IsDown && this.backend.CurrentFails > 0 {
			this.backend.CurrentFails = 0
		}
	}

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

	// 响应回调
	if this.responseCallback != nil {
		this.responseCallback(writer)
	}

	// 准备
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	// 分析API中的status
	if this.api != nil {
		statusCode := resp.Header.Get("Tea-Status-Code")
		if len(statusCode) == 0 && this.server.API.StatusScriptOn && len(this.server.API.StatusScript) > 0 {
			data, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				statusCode, _ = StatusCodeParser(resp.StatusCode, writer.Header(), data, this.server.API.StatusScript)
				resp.Body = ioutil.NopCloser(bytes.NewReader(data))
			}
		}
		if len(statusCode) > 0 {
			this.responseAPIStatus = statusCode
		}
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
		this.addError(err)
		return nil
	}
	return nil
}

// 调用代理
func (this *Request) callProxy(writer *ResponseWriter) error {
	options := maps.Map{
		"request":   this.raw,
		"formatter": this.Format,
	}
	backend := this.proxy.NextBackend(options)

	responseCallback := options.Get("responseCallback")
	if responseCallback != nil {
		f, ok := responseCallback.(func(http.ResponseWriter))
		if ok {
			this.responseCallback = f
		}
	}

	this.backend = backend
	return this.callBackend(writer)
}

// 调用Fastcgi
func (this *Request) callFastcgi(writer *ResponseWriter) error {
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
		this.addError(err)
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
	fcgiReq.SetTimeout(this.fastcgi.ReadTimeoutDuration())
	fcgiReq.SetParams(params)
	fcgiReq.SetBody(this.raw.Body, uint32(this.requestLength()))

	resp, stderr, err := client.Call(fcgiReq)
	if err != nil {
		this.serverError(writer)
		//if this.debug {
		logs.Error(err)
		this.addError(err)
		//}
		return nil
	}

	if len(stderr) > 0 {
		logs.Println("Fastcgi Error: " + string(stderr))
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

	// 准备
	writer.Prepare(resp.ContentLength)

	// 设置响应码
	writer.WriteHeader(resp.StatusCode)

	// 输出内容
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
		this.addError(err)
		return nil
	}

	return nil
}

// 调用Rewrite
func (this *Request) callRewrite(writer *ResponseWriter) error {
	query := this.requestQueryString()
	target := this.rewriteReplace
	if len(query) > 0 {
		if strings.Index(target, "?") > 0 {
			target += "&" + query
		} else {
			target += "?" + query
		}
	}

	if this.rewriteRedirectMode == teaconfigs.RewriteFlagRedirect {
		// 跳转
		http.Redirect(writer, this.raw, target, http.StatusTemporaryRedirect)
		return nil
	}

	if this.rewriteRedirectMode == teaconfigs.RewriteFlagProxy {
		req, err := http.NewRequest(this.requestMethod(), target, this.raw.Body)
		if err != nil {
			return err
		}

		// ip
		remoteAddr := this.requestRemoteAddr()
		if len(remoteAddr) > 0 {
			index := strings.Index(this.raw.RemoteAddr, ":")
			ip := ""
			if index > -1 {
				ip = this.raw.RemoteAddr[:index]
			} else {
				ip = this.raw.RemoteAddr
			}
			req.Header.Set("X-Real-IP", ip)
			req.Header.Set("X-Forwarded-For", ip)
			req.Header.Set("X-Forwarded-By", ip)
		}

		// headers
		for _, h := range this.headers {
			req.Header.Add(h.Name, h.Value)
		}

		var client *http.Client = nil
		if len(req.Host) > 0 {
			host := req.Host
			if !strings.Contains(host, ":") {
				if req.URL.Scheme == "https" {
					host += ":443"
				} else {
					host += ":80"
				}
			}
			client = SharedClientPool.client("", host, 30*time.Second, 0, 0)
		} else {
			client = &http.Client{
				Timeout: 30 * time.Second,
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			logs.Error(errors.New(req.URL.String() + ": " + err.Error()))
			this.addError(err)
			this.serverError(writer)
			return err
		}
		defer resp.Body.Close()

		// Header
		writer.AddHeaders(resp.Header)
		writer.Prepare(resp.ContentLength)

		// 设置响应代码
		writer.WriteHeader(resp.StatusCode)

		// 输出内容
		_, err = io.Copy(writer, resp.Body)

		return err
	}

	return nil
}

// 调用API Mock
func (this *Request) callMock(writer *ResponseWriter) error {
	if this.api != nil && len(this.api.MockFiles) > 0 {
		mock := this.api.RandMock()
		if mock != nil {
			for _, header := range mock.Headers {
				name := header.GetString("name")
				value := header.GetString("value")
				if len(name) > 0 {
					writer.Header().Set(name, value)
				}
			}

			writer.Header().Set("Tea-API-Mock", "on")

			if len(mock.File) > 0 {
				reader, err := files.NewReader(Tea.ConfigFile(mock.File))
				if err == nil {
					defer reader.Close()
					data := reader.ReadAll()
					writer.Write(data)
				}
			} else {
				writer.Write([]byte(mock.Text))
			}
		}
	} else {
		writer.Write([]byte("mock data not found"))
	}

	return nil
}

// 处理API
func (this *Request) consumeAPI(writer *ResponseWriter) bool {
	if len(this.api.AuthType) == 0 {
		this.api.AuthType = apiconfig.APIAuthTypeNone
	}
	consumer, authorized := this.server.API.FindConsumerForRequest(this.api.AuthType, this.raw)
	if !authorized {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unauthorized Request"))
		return false
	}
	if consumer == nil {
		return true
	}

	if !consumer.AllowAPI(this.api.Path) {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden Request"))
		return false
	}

	if !consumer.Policy.AllowAccess(this.requestRemoteAddr()) {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden Request"))
		return false
	}

	reason, allowed := consumer.Policy.AllowTraffic()
	if !allowed {
		writer.WriteHeader(http.StatusTooManyRequests)
		writer.Write([]byte("[" + reason + "]Request Quota Exceeded"))
		return false
	}

	return true
}

func (this *Request) notFoundError(writer *ResponseWriter) {
	msg := "404 page not found: '" + this.requestURI() + "'"

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(msg))

}

func (this *Request) serverError(writer *ResponseWriter) {
	statusCode := http.StatusInternalServerError

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 自定义Header
	for _, header := range this.headers {
		if header.Match(statusCode) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			writer.Header().Set(header.Name, header.Value)
		}
	}

	writer.WriteHeader(statusCode)
	writer.Write([]byte(http.StatusText(statusCode)))

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

func (this *Request) allowGzip() bool {
	encodingList := this.raw.Header.Get("Accept-Encoding")
	if len(encodingList) == 0 {
		return false
	}
	encodings := strings.Split(encodingList, ",")
	for _, encoding := range encodings {
		if encoding == "gzip" {
			return true
		}
	}
	return false
}

func (this *Request) CachePolicy() *shared.CachePolicy {
	return this.cachePolicy
}

func (this *Request) SetCachePolicy(config *shared.CachePolicy) {
	this.cachePolicy = config
}

func (this *Request) SetCacheEnabled() {
	this.cacheEnabled = true
}

// 判断缓存策略是否有效
func (this *Request) IsCacheEnabled() bool {
	return this.cacheEnabled
}

// 设置监控状态
func (this *Request) SetIsWatching(isWatching bool) {
	this.isWatching = isWatching
}

// 判断是否在监控
func (this *Request) IsWatching() bool {
	return this.isWatching
}

// 利用请求参数格式化字符串
func (this *Request) Format(source string) string {
	if len(source) == 0 {
		return ""
	}

	var hasVarMapping = len(this.varMapping) > 0

	return teautils.ParseVariables(source, func(varName string) string {
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
			return fmt.Sprintf("%d", this.responseWriter.SentBodyBytes()) // TODO 加上Header长度
		case "bodyBytesSent":
			return fmt.Sprintf("%d", this.responseWriter.SentBodyBytes())
		case "status":
			return fmt.Sprintf("%d", this.responseWriter.StatusCode())
		case "statusMessage":
			return http.StatusText(this.responseWriter.StatusCode())
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
			return "${" + varName + "}"
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

		// backend.
		if prefix == "backend" {
			if this.backend != nil {
				switch suffix {
				case "address":
					return this.backend.Address
				case "id":
					return this.backend.Id
				case "scheme":
					return this.backend.Scheme
				case "code":
					return this.backend.Code
				}
			}
			return ""
		}

		return "${" + varName + "}"
	})
}

// 格式化一组字符串
func (this *Request) formatAll(sources []string) []string {
	result := []string{}
	for _, s := range sources {
		result = append(result, this.Format(s))
	}
	return result
}

// 记录日志
func (this *Request) log() {
	// 计算请求时间
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	if !this.enableAccessLog {
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
		BytesSent:       this.responseWriter.SentBodyBytes(), // TODO 加上Header Size
		BodyBytesSent:   this.responseWriter.SentBodyBytes(),
		Status:          this.responseWriter.StatusCode(),
		StatusMessage:   "",
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
		Errors:          this.errors,
		HasErrors:       len(this.errors) > 0,
	}
	accessLog.SetShouldStat(true)

	if this.api != nil {
		accessLog.APIPath = this.api.Path
		accessLog.APIStatus = this.responseAPIStatus
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

	accessLog.SentHeader = this.responseWriter.Header()

	if len(this.requestData) > 0 {
		accessLog.RequestData = this.requestData
	}

	if this.responseWriter.BodyIsCopying() {
		accessLog.ResponseHeaderData = this.responseWriter.HeaderData()
		accessLog.ResponseBodyData = this.responseWriter.Body()
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
			indexFiles, err := filepath.Glob(dir + Tea.DS + index)
			if err != nil {
				logs.Error(err)
				this.addError(err)
				continue
			}
			if len(indexFiles) > 0 {
				return filepath.Base(indexFiles[0])
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

func (this *Request) addVarMapping(varMapping map[string]string) {
	for k, v := range varMapping {
		this.varMapping[k] = v
	}
}

func (this *Request) addError(err error) {
	if err == nil {
		return
	}
	this.errors = append(this.errors, err.Error())
}
