package teaproxy

import (
	"crypto/tls"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Scheme = uint8

const (
	SchemeHTTP  = Scheme(1)
	SchemeHTTPS = Scheme(2)
)

// 代理服务监听器
type Listener struct {
	httpServer *http.Server

	IsChanged bool // 标记是否改变，用来在其他地方重启改变的监听器

	Scheme  Scheme // http & https
	Address string
	Error   error

	servers        []*teaconfigs.ServerConfig // 待启用的server
	currentServers []*teaconfigs.ServerConfig // 当前可用的Server
	namedServers   map[string]*NamedServer    // 域名 => server

	serversLocker      sync.RWMutex
	namedServersLocker sync.RWMutex
}

// 获取新对象
func NewListener() *Listener {
	return &Listener{
		namedServers: map[string]*NamedServer{},
	}
}

// 应用配置
func (this *Listener) ApplyServer(server *teaconfigs.ServerConfig) {
	this.serversLocker.Lock()
	defer this.serversLocker.Unlock()

	this.IsChanged = true

	isAvailable := false
	if this.Scheme == SchemeHTTP && server.Http {
		isAvailable = true
	} else if this.Scheme == SchemeHTTPS && server.SSL != nil && server.SSL.On {
		isAvailable = true
	}

	if !isAvailable {
		// 删除
		result := []*teaconfigs.ServerConfig{}
		for _, s := range this.servers {
			if s.Id == server.Id {
				continue
			}
			result = append(result, s)
		}
		this.servers = result

		return
	}

	found := false
	for index, s := range this.servers {
		if s.Id == server.Id {
			this.servers[index] = server
			found = true
			break
		}
	}
	if !found {
		this.servers = append(this.servers, server)
	}
}

// 删除配置
func (this *Listener) RemoveServer(serverId string) {
	this.serversLocker.Lock()
	defer this.serversLocker.Unlock()

	this.IsChanged = true
	result := []*teaconfigs.ServerConfig{}
	for _, s := range this.servers {
		if s.Id == serverId {
			continue
		}
		result = append(result, s)
	}
	this.servers = result
}

// 重置所有配置
func (this *Listener) Reset() {
	this.serversLocker.Lock()
	defer this.serversLocker.Unlock()

	this.IsChanged = true
	this.servers = []*teaconfigs.ServerConfig{}
}

// 判断是否包含某个配置
func (this *Listener) HasServer(serverId string) bool {
	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	for _, s := range this.servers {
		if s.Id == serverId {
			return true
		}
	}
	return false
}

// 是否包含配置
func (this *Listener) HasServers() bool {
	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	return len(this.servers) > 0
}

// 启动
func (this *Listener) Start() error {
	return this.Reload()
}

// 刷新
func (this *Listener) Reload() error {
	this.namedServersLocker.Lock()
	this.namedServers = map[string]*NamedServer{}
	this.namedServersLocker.Unlock()

	this.serversLocker.Lock()
	this.currentServers = this.servers
	hasServers := len(this.currentServers) > 0
	this.IsChanged = false
	this.Error = nil

	if !hasServers {
		defer this.serversLocker.Unlock()

		// 检查是否已启动
		if this.httpServer != nil {
			return this.Shutdown()
		}

		return nil
	} else {
		this.serversLocker.Unlock()
	}

	// 如果已经启动，则不做任何事情
	if this.httpServer != nil {
		return nil
	}

	// 如果没启动，则启动
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// QPS计算
		atomic.AddInt32(&qps, 1)

		// 处理
		this.handle(writer, req)
	})

	var err error

	this.httpServer = &http.Server{
		Addr:        this.Address,
		Handler:     httpHandler,
		IdleTimeout: 2 * time.Minute,
	}
	this.httpServer.SetKeepAlivesEnabled(true)

	if this.Scheme == SchemeHTTP {
		logs.Println("start listener on", this.Address)
		err = this.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.Error(errors.New("[listener]" + this.Address + ": " + err.Error()))
		} else {
			err = nil
		}
	}

	if this.Scheme == SchemeHTTPS {
		logs.Println("start ssl listener on", this.Address)

		this.httpServer.TLSConfig = &tls.Config{
			Certificates: nil,
			GetConfigForClient: func(info *tls.ClientHelloInfo) (config *tls.Config, e error) {
				ssl, _, err := this.matchSSL(info.ServerName)
				if err != nil {
					return nil, err
				}

				cipherSuites := ssl.TLSCipherSuites()
				if len(cipherSuites) == 0 {
					cipherSuites = nil
				}
				return &tls.Config{
					Certificates: nil,
					MinVersion:   ssl.TLSMinVersion(),
					CipherSuites: cipherSuites,
					GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
						_, cert, err := this.matchSSL(info.ServerName)
						if err != nil {
							return nil, err
						}
						if cert == nil {
							return nil, errors.New("[listener]no certs found for '" + info.ServerName + "'")
						}
						return cert, nil
					},
					NextProtos: []string{http2.NextProtoTLS},
				}, nil
			},
			GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
				_, cert, err := this.matchSSL(info.ServerName)
				if err != nil {
					return nil, err
				}
				if cert == nil {
					return nil, errors.New("[listener]no certs found for '" + info.ServerName + "'")
				}
				return cert, nil
			},
		}

		// support http/2
		http2.ConfigureServer(this.httpServer, nil)

		err = this.httpServer.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			logs.Error(errors.New("[listener]" + this.Address + ": " + err.Error()))
		} else {
			err = nil
		}
	}

	this.httpServer = nil

	return err
}

// 关闭
func (this *Listener) Shutdown() error {
	if this.httpServer != nil {
		logs.Println("shutdown listener on", this.Address)
		return this.httpServer.Shutdown(context.Background())
	}
	return nil
}

// 处理请求
func (this *Listener) handle(writer http.ResponseWriter, rawRequest *http.Request) {
	responseWriter := NewResponseWriter(writer)

	// 插件过滤
	if teaplugins.HasRequestFilters {
		result, willContinue := teaplugins.FilterRequest(rawRequest)
		if !willContinue {
			return
		}
		rawRequest = result
	}

	// 域名
	reqHost := rawRequest.Host
	domain, _, err := net.SplitHostPort(reqHost)
	if err != nil {
		domain = reqHost
	}
	server, serverName := this.findNamedServer(domain)
	if server == nil {
		http.Error(writer, "404 page not found: '"+rawRequest.URL.String()+"'", http.StatusNotFound)
		return
	}

	// 包装新的请求
	req := NewRequest(rawRequest)
	req.host = reqHost
	req.method = rawRequest.Method
	req.uri = rawRequest.URL.RequestURI()
	if this.Scheme == SchemeHTTP {
		req.rawScheme = "http"
	} else if this.Scheme == SchemeHTTPS {
		req.rawScheme = "https"
	} else {
		req.rawScheme = "http"
	}
	req.scheme = "http" // 转发后的scheme
	req.serverName = serverName
	req.serverAddr = this.Address
	req.root = server.Root
	req.index = server.Index
	req.charset = server.Charset

	// 配置请求
	err = req.configure(server, 0)
	if err != nil {
		req.serverError(responseWriter)
		logs.Error(errors.New(reqHost + rawRequest.URL.String() + ": " + err.Error()))
		return
	}

	// 处理请求
	req.call(responseWriter)
}

// 根据域名来查找匹配的域名
func (this *Listener) findNamedServer(name string) (serverConfig *teaconfigs.ServerConfig, serverName string) {
	// 读取缓存
	this.namedServersLocker.RLock()
	namedServer, found := this.namedServers[name]
	if found {
		this.namedServersLocker.RUnlock()
		return namedServer.Server, namedServer.Name
	}
	this.namedServersLocker.RUnlock()

	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	countServers := len(this.currentServers)
	if countServers == 0 {
		return nil, ""
	}

	// 只记录N个记录，防止内存耗尽
	maxNamedServers := 10240

	// 如果只有一个server，则默认为这个
	if countServers == 1 {
		server := this.currentServers[0]
		matchedName, matched := server.MatchName(name)
		if matched {
			if len(matchedName) > 0 {
				this.namedServersLocker.Lock()
				if len(this.namedServers) < maxNamedServers {
					this.namedServers[name] = &NamedServer{
						Name:   matchedName,
						Server: server,
					}
				}
				this.namedServersLocker.Unlock()
				return server, matchedName
			} else {
				return server, name
			}
		}

		// 匹配第一个域名
		firstName := server.FirstName()
		if len(firstName) > 0 {
			return server, firstName
		}
		return server, name
	}

	// 精确查找
	for _, server := range this.currentServers {
		if lists.ContainsString(server.Name, name) {
			this.namedServersLocker.Lock()
			if len(this.namedServers) < maxNamedServers {
				this.namedServers[name] = &NamedServer{
					Name:   name,
					Server: server,
				}
			}
			this.namedServersLocker.Unlock()
			return server, name
		}
	}

	// 模糊查找
	for _, server := range this.currentServers {
		if _, matched := server.MatchName(name); matched {
			this.namedServersLocker.Lock()
			if len(this.namedServers) < maxNamedServers {
				this.namedServers[name] = &NamedServer{
					Name:   name,
					Server: server,
				}
			}
			this.namedServersLocker.Unlock()
			return server, name
		}
	}

	// 如果没有找到，则匹配到第一个
	server := this.currentServers[0]
	firstName := server.FirstName()
	if len(firstName) > 0 {
		this.namedServersLocker.Lock()
		if len(this.namedServers) < maxNamedServers {
			this.namedServers[name] = &NamedServer{
				Name:   firstName,
				Server: server,
			}
		}
		this.namedServersLocker.Unlock()
		return server, firstName
	}

	return server, name
}

// 根据域名匹配证书
func (this *Listener) matchSSL(domain string) (*teaconfigs.SSLConfig, *tls.Certificate, error) {
	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	if len(domain) == 0 {
		if len(this.currentServers) > 0 && this.currentServers[0].SSL != nil {
			logs.Error(errors.New("[listener]no tls server name found"))
			return this.currentServers[0].SSL, this.currentServers[0].SSL.FirstCert(), nil
		}
		return nil, nil, errors.New("[listener]no tls server name found")
	}

	// 通过代理服务域名配置匹配
	server, _ := this.findNamedServer(domain)
	if server == nil || server.SSL == nil || !server.SSL.On {
		// 搜索所有的Server，通过SSL证书内容中的DNSName匹配
		for _, server := range this.currentServers {
			if server.SSL == nil || !server.SSL.On {
				continue
			}
			cert, ok := server.SSL.MatchDomain(domain)
			if ok {
				return server.SSL, cert, nil
			}
		}

		return nil, nil, errors.New("[listener]no server found for '" + domain + "'")
	}

	// 证书是否匹配
	cert, ok := server.SSL.MatchDomain(domain)
	if ok {
		return server.SSL, cert, nil
	}

	return server.SSL, server.SSL.FirstCert(), nil
}
