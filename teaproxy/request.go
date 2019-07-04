package teaproxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teawaf"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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
// HTTP HEADER RFC: https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html
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
	requestHeaders  []*shared.HeaderConfig // 自定义请求Header
	responseHeaders []*shared.HeaderConfig // 自定义响应Header
	ignoreHeaders   []string               // 忽略的响应Header
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

	waf *teawaf.WAF

	pages          []*teaconfigs.PageConfig
	shutdownPageOn bool
	shutdownPage   string

	rewriteId           string // 匹配的rewrite id
	rewriteReplace      string // 经过rewrite之后的URL
	rewriteRedirectMode string // 跳转方式
	rewriteIsExternal   bool   // 是否为外部URL

	redirectToHttps bool

	websocket *teaconfigs.WebsocketConfig

	tunnel *teaconfigs.TunnelConfig

	// 执行请求
	filePath string

	responseWriter   *ResponseWriter
	responseCallback func(http.ResponseWriter)

	requestFromTime time.Time // 请求开始时间
	requestCost     float64   // 请求耗时
	requestMaxSize  int64

	isWatching  bool     // 是否在监控
	requestData []byte   // 导出的request，在监控请求的时候有用
	errors      []string // 错误信息

	enableStat bool
	accessLog  *teaconfigs.AccessLogConfig

	gzipLevel     uint8
	gzipMinLength int64
	debug         bool

	hasForwardHeader bool
}

// 获取新的请求
func NewRequest(rawRequest *http.Request) *Request {
	now := time.Now()

	req := &Request{
		varMapping:      map[string]string{},
		raw:             rawRequest,
		rawURI:          rawRequest.URL.RequestURI(),
		requestFromTime: now,
		enableStat:      true,
		attrs:           map[string]string{},
	}

	backendCall := shared.NewRequestCall()
	backendCall.Request = rawRequest
	backendCall.Formatter = req.Format
	req.backendCall = backendCall
	_, req.hasForwardHeader = rawRequest.Header["X-Forwarded-For"]

	return req
}

func (this *Request) configure(server *teaconfigs.ServerConfig, redirects int) error {
	isChanged := this.server != server
	this.server = server

	if redirects > 8 {
		return errors.New("too many redirects")
	}
	redirects++

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
	if server.HasRequestHeaders() {
		this.requestHeaders = append(this.requestHeaders, server.RequestHeaders...)
	}

	if server.HasResponseHeaders() {
		this.responseHeaders = append(this.responseHeaders, server.Headers ...)
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

	// waf
	if server.WAFOn {
		waf := server.WAF()
		if waf != nil && waf.On {
			this.waf = waf
		}
	} else {
		this.waf = nil
	}

	// tunnel
	if server.Tunnel != nil && server.Tunnel.On {
		this.tunnel = server.Tunnel
		return nil
	} else {
		this.tunnel = nil
	}

	// other
	if server.MaxBodyBytes() > 0 {
		this.requestMaxSize = server.MaxBodyBytes()
	}
	if len(server.AccessLog) > 0 {
		this.accessLog = server.AccessLog[0]
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

	if server.RedirectToHttps && this.rawScheme == "http" {
		this.redirectToHttps = true
		return nil
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
			if len(location.AccessLog) > 0 {
				this.accessLog = location.AccessLog[0]
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
			if location.RedirectToHttps && this.rawScheme == "http" {
				this.redirectToHttps = true
				return nil
			}

			if location.CacheOn {
				cachePolicy := location.CachePolicyObject()
				if cachePolicy != nil && cachePolicy.On {
					this.cachePolicy = cachePolicy
				}
			} else {
				this.cachePolicy = nil
			}

			if location.WAFOn {
				waf := location.WAF()
				if waf != nil && waf.On {
					this.waf = waf
				}
			} else {
				this.waf = nil
			}

			if location.HasRequestHeaders() {
				this.requestHeaders = append(this.requestHeaders, location.RequestHeaders...)
			}

			if location.HasResponseHeaders() {
				this.responseHeaders = append(this.responseHeaders, location.Headers ...)
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

						if rule.HasResponseHeaders() {
							this.responseHeaders = append(this.responseHeaders, rule.Headers...)
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

				if fastcgi.HasResponseHeaders() {
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

				if backend.HasRequestHeaders() {
					this.requestHeaders = append(this.requestHeaders, backend.RequestHeaders...)
				}

				if backend.HasResponseHeaders() {
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

				if rule.HasRequestHeaders() {
					this.requestHeaders = append(this.requestHeaders, rule.RequestHeaders...)
				}

				if rule.HasResponseHeaders() {
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

		if fastcgi.HasRequestHeaders() {
			this.requestHeaders = append(this.requestHeaders, fastcgi.RequestHeaders...)
		}

		if fastcgi.HasResponseHeaders() {
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
		if backend.HasRequestHeaders() {
			this.requestHeaders = append(this.requestHeaders, backend.RequestHeaders...)
		}

		if backend.HasResponseHeaders() {
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

	// WAF
	if this.waf != nil {
		if this.callWAFRequest(writer) {
			return nil
		}
	}

	// 跳转到https
	if this.redirectToHttps {
		this.callRedirectToHttps(writer)
		return nil
	}

	// 临时关闭页面
	if this.shutdownPageOn {
		return this.callShutdown(writer)
	}

	// hook
	b := CallRequestBeforeHook(this, writer)
	if !b {
		return nil
	}

	// gzip压缩
	if this.gzipLevel > 0 && this.allowGzip() {
		writer.Gzip(this.gzipLevel, this.gzipMinLength)
		defer writer.Close()
	}

	// watch
	if this.isWatching {
		// 判断如果Content-Length过长，则截断
		reqData, err := httputil.DumpRequest(this.raw, true)
		if err == nil {
			if len(reqData) > 100240 {
				reqData = reqData[:100240]
			}
			this.requestData = reqData
		}

		writer.SetBodyCopying(true)
	} else {
		max := 512 * 1024 // 512K
		if this.accessLog != nil && lists.ContainsInt(this.accessLog.Fields, tealogs.AccessLogFieldRequestBody) {
			body, err := ioutil.ReadAll(this.raw.Body)
			if err == nil {
				if len(body) > max {
					this.requestData = body[:max]
				} else {
					this.requestData = body
				}
			}
			this.raw.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		if this.accessLog != nil && lists.ContainsInt(this.accessLog.Fields, tealogs.AccessLogFieldResponseBody) {
			writer.SetBodyCopying(true)
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

	if this.tunnel != nil {
		return this.callTunnel(writer)
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
			writer.Header().Set(header.Name, this.Format(header.Value))
		}
	}

	writer.WriteHeader(statusCode)
	writer.Write([]byte(http.StatusText(statusCode)))
}

func (this *Request) requestRemoteAddr() string {
	// X-Forwarded-For
	forwardedFor := this.raw.Header.Get("X-Forwarded-For")
	if len(forwardedFor) > 0 {
		commaIndex := strings.Index(forwardedFor, ",")
		if commaIndex > 0 {
			return forwardedFor[:commaIndex]
		}
		return forwardedFor
	}

	// Real-IP
	{
		realIP, ok := this.raw.Header["X-Real-IP"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Real-Ip
	{
		realIP, ok := this.raw.Header["X-Real-Ip"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Remote-Addr
	remoteAddr := this.raw.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	} else {
		return remoteAddr
	}
}

func (this *Request) requestRemotePort() int {
	_, port, err := net.SplitHostPort(this.raw.RemoteAddr)
	if err == nil {
		types.Int(port)
	}
	return 0
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
	_, port, err := net.SplitHostPort(this.serverAddr)
	if err == nil {
		return types.Int(port)
	}
	return 0
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

// 输出自定义Response Header
func (this *Request) WriteResponseHeaders(writer *ResponseWriter, statusCode int) {
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	responseHeader := writer.Header()

	for _, header := range this.responseHeaders {
		if !header.On {
			continue
		}
		if header.Match(statusCode) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			if header.HasVariables() {
				responseHeader.Set(header.Name, this.Format(header.Value))
			} else {
				responseHeader.Set(header.Name, header.Value)
			}
		}
	}

	// hsts
	if this.rawScheme == "https" &&
		this.server.SSL != nil &&
		this.server.SSL.On &&
		this.server.SSL.HSTS != nil &&
		this.server.SSL.HSTS.On &&
		this.server.SSL.HSTS.Match(this.host) {
		responseHeader.Set(this.server.SSL.HSTS.HeaderKey(), this.server.SSL.HSTS.HeaderValue())
	}
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
		case "rawRemoteAddr":
			addr := this.raw.RemoteAddr
			host, _, err := net.SplitHostPort(addr)
			if err == nil {
				addr = host
			}
			return addr
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

		// node
		if prefix == "node" {
			node := teaconfigs.SharedNodeConfig()
			if node != nil {
				switch suffix {
				case "id":
					return node.Id
				case "name":
					return node.Name
				case "role":
					return node.Role
				}
			}
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

	if (this.accessLog == nil || !this.accessLog.On) && !this.enableStat {
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
	accessLog.SetShouldWrite(this.accessLog != nil && this.accessLog.On && this.accessLog.Match(this.responseWriter.statusCode))
	accessLog.SetShouldStat(this.enableStat)
	if this.accessLog != nil {
		accessLog.SetWritingFields(this.accessLog.Fields)
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

// 添加自定义变量
func (this *Request) SetVarMapping(varName string, varValue string) {
	this.varMapping[varName] = varValue
}

func (this *Request) addError(err error) {
	if err == nil {
		return
	}
	this.errors = append(this.errors, err.Error())
}

// 设置代理相关头部信息
// 参考：https://tools.ietf.org/html/rfc7239
func (this *Request) setProxyHeaders(header http.Header) {
	delete(header, "Connection")

	remoteAddr := this.raw.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		remoteAddr = host
	}

	// x-real-ip
	{
		_, ok1 := header["X-Real-IP"]
		_, ok2 := header["X-Real-Ip"]
		if !ok1 && !ok2 {
			header["X-Real-IP"] = []string{remoteAddr}
		}
	}

	// X-Forwarded-For
	{
		forwardedFor, ok := header["X-Forwarded-For"]
		if ok {
			if this.hasForwardHeader {
				header["X-Forwarded-For"] = []string{strings.Join(forwardedFor, ", ") + ", " + remoteAddr}
			}
		} else {
			header["X-Forwarded-For"] = []string{remoteAddr}
		}
	}

	// Forwarded
	/**{
		forwarded, ok := header["Forwarded"]
		if ok {
			header["Forwarded"] = []string{strings.Join(forwarded, ", ") + ", by=" + this.serverAddr + "; for=" + remoteAddr + "; host=" + this.host + "; proto=" + this.rawScheme}
		} else {
			header["Forwarded"] = []string{"by=" + this.serverAddr + "; for=" + remoteAddr + "; host=" + this.host + "; proto=" + this.rawScheme}
		}
	}**/

	// others
	this.raw.Header.Set("X-Forwarded-By", this.serverAddr)

	if _, ok := header["X-Forwarded-Host"]; !ok {
		this.raw.Header.Set("X-Forwarded-Host", this.host)
	}

	if _, ok := header["X-Forwarded-Proto"]; !ok {
		this.raw.Header.Set("X-Forwarded-Proto", this.rawScheme)
	}
}
