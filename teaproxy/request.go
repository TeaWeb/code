package teaproxy

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
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

	attrs map[string]string // 附加参数

	scheme          string
	rawScheme       string // 原始的scheme
	uri             string
	rawURI          string // 跳转之前的uri
	host            string
	method          string
	serverName      string // @TODO
	serverAddr      string
	charset         string
	responseHeaders []*shared.HeaderConfig // 自定义响应Header
	ignoreHeaders   []string               // 忽略的Header
	varMapping      map[string]string      // 自定义变量

	root  string   // 资源根目录
	index []string // 目录下默认访问的文件

	backend     *teaconfigs.BackendConfig
	backendCall *shared.RequestCall

	fastcgi      *teaconfigs.FastcgiConfig
	proxy        *teaconfigs.ServerConfig
	location     *teaconfigs.LocationConfig
	accessPolicy *shared.AccessPolicy

	cachePolicy  *shared.CachePolicy
	cacheEnabled bool

	pages          []*teaconfigs.PageConfig
	shutdownPageOn bool
	shutdownPage   string

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

	requestFromTime time.Time // 请求开始时间
	requestCost     float64   // 请求耗时
	requestMaxSize  int64

	isWatching        bool     // 是否在监控
	requestData       []byte   // 导出的request，在监控请求的时候有用
	responseAPIStatus string   // API状态码
	errors            []string // 错误信息

	enableAccessLog bool
	enableStat      bool
	accessLogFields []int

	gzipLevel     uint8
	gzipMinLength int64
	debug         bool
}

// 获取新的请求
func NewRequest(rawRequest *http.Request) *Request {
	now := time.Now()

	req := &Request{
		varMapping:      map[string]string{},
		raw:             rawRequest,
		rawURI:          rawRequest.URL.RequestURI(),
		requestFromTime: now,
		enableAccessLog: true,
		enableStat:      true,
		attrs:           map[string]string{},
	}

	backendCall := shared.NewRequestCall()
	backendCall.Request = rawRequest
	backendCall.Formatter = req.Format
	req.backendCall = backendCall

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
	if server.HasHeaders() {
		// 延迟执行，让Header有机会加入Backend, Fastcgi等信息
		defer func() {
			this.responseHeaders = append(this.responseHeaders, server.FormatHeaders(this.Format) ...)
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
	} else {
		this.cachePolicy = nil
	}

	// other
	if server.MaxBodyBytes() > 0 {
		this.requestMaxSize = server.MaxBodyBytes()
	}
	if server.DisableAccessLog {
		this.enableAccessLog = false
	} else {
		this.accessLogFields = server.AccessLogFields
	}
	if server.DisableStat {
		this.enableStat = false
	}
	if len(server.Pages) > 0 {
		this.pages = append(this.pages, server.Pages ...)
	}
	if server.ShutdownPageOn {
		this.shutdownPageOn = true
		this.shutdownPage = server.ShutdownPage
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
			if api.HasHeaders() {
				this.responseHeaders = append(this.responseHeaders, api.FormatHeaders(func(source string) string {
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
		if locationMatches, ok := location.Match(path, this.Format); ok {
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
			} else {
				this.accessLogFields = location.AccessLogFields
			}
			if location.DisableStat {
				this.enableStat = false
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
			} else {
				this.cachePolicy = nil
			}

			if location.HasHeaders() {
				this.responseHeaders = append(this.responseHeaders, location.FormatHeaders(this.Format) ...)
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

						if rule.HasHeaders() {
							this.responseHeaders = append(this.responseHeaders, rule.FormatHeaders(func(source string) string {
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
							server := SharedManager.FindServer(proxyId)
							if server == nil {
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

				if fastcgi.HasHeaders() {
					this.responseHeaders = append(this.responseHeaders, fastcgi.Headers ...)
				}

				if len(fastcgi.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, fastcgi.IgnoreHeaders ...)
				}

				continue
			}

			// proxy
			if len(location.Proxy) > 0 {
				server := SharedManager.FindServer(location.Proxy)
				if server == nil {
					return errors.New("server with '" + location.Proxy + "' not found")
				}
				if !server.On {
					return errors.New("server with '" + location.Proxy + "' not available now")
				}
				return this.configure(server, redirects)
			}

			// backends
			if len(location.Backends) > 0 {
				backend := location.NextBackend(this.backendCall)
				if backend == nil {
					return errors.New("no backends available")
				}
				if len(this.backendCall.ResponseCallbacks) > 0 {
					this.responseCallback = this.backendCall.CallResponseCallbacks
				}
				this.backend = backend
				locationConfigured = true

				if backend.HasHeaders() {
					this.responseHeaders = append(this.responseHeaders, backend.Headers ...)
				}

				if len(backend.IgnoreHeaders) > 0 {
					this.ignoreHeaders = append(this.ignoreHeaders, backend.IgnoreHeaders ...)
				}

				continue
			}

			// websocket
			if location.Websocket != nil && location.Websocket.On {
				this.backend = location.Websocket.NextBackend(this.backendCall)
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

				if rule.HasHeaders() {
					this.responseHeaders = append(this.responseHeaders, rule.Headers ...)
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
					server := SharedManager.FindServer(proxyId)
					if server == nil {
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

		if fastcgi.HasHeaders() {
			this.responseHeaders = append(this.responseHeaders, fastcgi.Headers ...)
		}

		if len(fastcgi.IgnoreHeaders) > 0 {
			this.ignoreHeaders = append(this.ignoreHeaders, fastcgi.IgnoreHeaders ...)
		}

		return nil
	}

	// proxy
	if len(server.Proxy) > 0 {
		server := SharedManager.FindServer(server.Proxy)
		if server == nil {
			return errors.New("server with '" + server.Proxy + "' not found")
		}
		if !server.On {
			return errors.New("server with '" + server.Proxy + "' not available now")
		}
		return this.configure(server, redirects)
	}

	// 转发到后端
	backend := server.NextBackend(this.backendCall)
	if backend == nil {
		if len(this.root) == 0 {
			return errors.New("no backends available")
		}
	}
	if len(this.backendCall.ResponseCallbacks) > 0 {
		this.responseCallback = this.backendCall.CallResponseCallbacks
	}
	this.backend = backend

	if backend != nil {
		if backend.HasHeaders() {
			this.responseHeaders = append(this.responseHeaders, backend.Headers ...)
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

	// 临时关闭页面
	if this.shutdownPageOn {
		return this.callShutdown(writer)
	}

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
	if this.callPage(writer, http.StatusNotFound) {
		return
	}

	msg := "404 page not found: '" + this.requestURI() + "'"

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(msg))
}

func (this *Request) serverError(writer *ResponseWriter) {
	if this.callPage(writer, http.StatusInternalServerError) {
		return
	}

	statusCode := http.StatusInternalServerError

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 自定义Header
	for _, header := range this.responseHeaders {
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
	return cookie.Value
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

// 设置URI
func (this *Request) SetURI(uri string) {
	this.uri = uri
}

// 设置Host
func (this *Request) SetHost(host string) {
	this.host = host
}

// 获取原始的请求
func (this *Request) Raw() *http.Request {
	return this.raw
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
			return fmt.Sprintf("%.6f", this.requestCost)
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
			return this.requestFromTime.Format("2006-01-02T15:04:05.000Z07:00")
		case "timeLocal":
			return this.requestFromTime.Format("2/Jan/2006:15:04:05 -0700")
		case "msec":
			return fmt.Sprintf("%.6f", float64(this.requestFromTime.Unix())+float64(this.requestFromTime.Nanosecond())/1000000000)
		case "timestamp":
			return fmt.Sprintf("%d", this.requestFromTime.Unix())
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

// 设置属性
func (this *Request) SetAttr(key string, value string) {
	this.attrs[key] = value
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
	this.requestCost = time.Since(this.requestFromTime).Seconds()

	if !this.enableAccessLog && !this.enableStat {
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
		RequestTime:     this.requestCost,
		RequestMethod:   this.requestMethod(),
		RequestFilename: this.requestFilename(),
		Scheme:          this.rawScheme,
		Proto:           this.requestProto(),
		BytesSent:       this.responseWriter.SentBodyBytes(), // TODO 加上Header Size
		BodyBytesSent:   this.responseWriter.SentBodyBytes(),
		Status:          this.responseWriter.StatusCode(),
		StatusMessage:   "",
		TimeISO8601:     this.requestFromTime.Format("2006-01-02T15:04:05.000Z07:00"),
		TimeLocal:       this.requestFromTime.Format("2/Jan/2006:15:04:05 -0700"),
		Msec:            float64(this.requestFromTime.Unix()) + float64(this.requestFromTime.Nanosecond())/1000000000,
		Timestamp:       this.requestFromTime.Unix(),
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
		Extend:          &tealogs.AccessLogExtend{},
		Attrs:           this.attrs,
	}

	// 日志和统计
	accessLog.SetShouldWrite(this.enableAccessLog)
	accessLog.SetShouldStat(this.enableStat)
	if this.enableAccessLog {
		accessLog.SetWritingFields(this.accessLogFields)
	}

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
