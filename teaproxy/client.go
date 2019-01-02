package teaproxy

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

// 客户端池单例
var SharedClientPool = NewClientPool()

// 客户端池
type ClientPool struct {
	clientsMap map[string]*http.Client // address => client
	locker     sync.Mutex
}

// 获取新对象
func NewClientPool() *ClientPool {
	return &ClientPool{
		clientsMap: map[string]*http.Client{},
	}
}

// 根据地址获取客户端
func (this *ClientPool) client(address string, connectionTimeout time.Duration, maxConnections uint) *http.Client {
	this.locker.Lock()
	defer this.locker.Unlock()

	client, found := this.clientsMap[address]
	if found {
		return client
	}

	// 超时时间
	if connectionTimeout <= 0 {
		connectionTimeout = 15 * time.Second
	}

	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 握手配置
			return (&net.Dialer{
				Timeout:   connectionTimeout,
				KeepAlive: 120 * time.Second,
				DualStack: true,
			}).DialContext(ctx, network, address)
		},
		MaxIdleConns:          int(maxConnections), // 0表示不限
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       0, // 不限
		TLSHandshakeTimeout:   0, // 不限
		ExpectContinueTimeout: 1 * time.Second,
	}

	c := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return &RedirectError{}
		},
	}
	this.clientsMap[address] = c

	return c
}
