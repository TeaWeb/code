package teaproxy

import (
	"context"
	"crypto/tls"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 客户端池单例
var SharedClientPool = NewClientPool()

// 客户端池
type ClientPool struct {
	clientsMap map[string]*http.Client // backend key => client
	locker     sync.RWMutex
}

// 获取新对象
func NewClientPool() *ClientPool {
	return &ClientPool{
		clientsMap: map[string]*http.Client{},
	}
}

// 根据地址获取客户端
func (this *ClientPool) client(backend *teaconfigs.BackendConfig) *http.Client {
	key := backend.UniqueKey()

	this.locker.RLock()
	client, found := this.clientsMap[key]
	if found {
		this.locker.RUnlock()
		return client
	}
	this.locker.RUnlock()
	this.locker.Lock()

	maxConnections := int(backend.MaxConns)
	connectionTimeout := backend.FailTimeoutDuration()
	address := backend.Address
	readTimeout := backend.ReadTimeoutDuration()
	idleTimeout := backend.IdleTimeoutDuration()
	idleConns := int(backend.IdleConns)

	// 超时时间
	if connectionTimeout <= 0 {
		connectionTimeout = 15 * time.Second
	}

	if idleTimeout <= 0 {
		idleTimeout = 2 * time.Minute
	}

	numberCPU := runtime.NumCPU()
	if numberCPU == 0 {
		numberCPU = 1
	}
	if maxConnections <= 0 {
		maxConnections = numberCPU
	}

	if idleConns <= 0 {
		idleConns = numberCPU
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 握手配置
			return (&net.Dialer{
				Timeout:   connectionTimeout,
				KeepAlive: 2 * time.Minute,
			}).DialContext(ctx, network, address)
		},
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   idleConns,
		MaxConnsPerHost:       maxConnections,
		IdleConnTimeout:       idleTimeout,
		ExpectContinueTimeout: 1 * time.Second,
		TLSHandshakeTimeout:   0, // 不限
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: nil,
	}

	client = &http.Client{
		Timeout:   readTimeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	this.clientsMap[key] = client

	// 关闭老的
	this.closeOldClient(key)

	this.locker.Unlock()

	return client
}

// 关闭老的client
func (this *ClientPool) closeOldClient(key string) {
	backendId := strings.Split(key, "@")[0]
	for key2, client := range this.clientsMap {
		backendId2 := strings.Split(key2, "@")[0]
		if backendId == backendId2 && key != key2 {
			teautils.CloseHTTPClient(client)
			delete(this.clientsMap, key2)
			break
		}
	}
}
