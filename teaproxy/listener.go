package teaproxy

import (
	"context"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"strings"
)

// 监听服务定义
type Listener struct {
	config *teaconfigs.ListenerConfig
	server *http.Server
	scheme string
}

// 新监听服务
func NewListener(config *teaconfigs.ListenerConfig) *Listener {
	listener := &Listener{
		config: config,
	}
	LISTENERS = append(LISTENERS, listener)
	return listener
}

// 启动
func (this *Listener) Start() {
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		this.handle(writer, req)
	})

	var err error

	this.server = &http.Server{
		Addr:    this.config.Address,
		Handler: httpHandler,
	}
	if this.config.SSL != nil && this.config.SSL.On {
		logs.Println("start ssl listener on", this.config.Address)
		this.scheme = "https"
		err = this.server.ListenAndServeTLS(Tea.ConfigFile(this.config.SSL.Certificate), Tea.ConfigFile(this.config.SSL.CertificateKey))
	}

	if this.config.Http {
		logs.Println("start listener on", this.config.Address)
		this.scheme = "http"
		err = this.server.ListenAndServe()
	}

	if err != nil {
		logs.Error(err)
		return
	}
}

// 关闭
func (this *Listener) Shutdown() error {
	if this.server != nil {
		return this.server.Shutdown(context.Background())
	}
	return nil
}

// 处理请求
func (this *Listener) handle(writer http.ResponseWriter, rawRequest *http.Request) {
	// 插件过滤
	result := teaplugins.FilterRequest(rawRequest)
	if !result {
		return
	}

	// 域名
	reqHost := rawRequest.Host
	colonIndex := strings.Index(reqHost, ":")
	domain := ""
	if colonIndex < 0 {
		domain = reqHost
	} else {
		domain = reqHost[:colonIndex]
	}
	server, serverName := this.config.FindNamedServer(domain)
	if server == nil {
		http.Error(writer, "404 page not found: '"+rawRequest.URL.String()+"'", http.StatusNotFound)
		return
	}

	// 包装新的请求
	req := NewRequest(rawRequest)
	req.host = reqHost
	req.method = rawRequest.Method
	req.uri = rawRequest.URL.RequestURI()
	req.rawScheme = this.scheme
	req.scheme = "http" // @TODO 支持 https
	req.serverName = serverName
	req.serverAddr = this.config.Address
	req.root = server.Root
	req.index = server.Index
	req.charset = server.Charset

	// 查找Location
	err := req.configure(server, 0)
	if err != nil {
		req.serverError(writer)
		logs.Error(errors.New(reqHost + rawRequest.URL.String() + ": " + err.Error()))
		return
	}

	// 处理请求
	req.Call(writer)
}
