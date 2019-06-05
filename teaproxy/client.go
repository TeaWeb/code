package teaproxy

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

// 客户端池单例
var SharedClientPool = NewClientPool()

// 客户端池
type ClientPool struct {
	clientsMap map[string]*http.Client // backend id => client
	locker     sync.RWMutex
}

// 获取新对象
func NewClientPool() *ClientPool {
	return &ClientPool{
		clientsMap: map[string]*http.Client{},
	}
}

// 根据地址获取客户端
func (this *ClientPool) client(backendId string, address string, connectionTimeout time.Duration, readTimeout time.Duration, maxConnections int32) *http.Client {
	key := backendId + "_" + address

	this.locker.RLock()
	client, found := this.clientsMap[key]
	if found {
		defer this.locker.RUnlock()
		return client
	}
	this.locker.RUnlock()

	// 超时时间
	if connectionTimeout <= 0 {
		connectionTimeout = 15 * time.Second
	}

	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 握手配置
			return (&net.Dialer{
				Timeout:   connectionTimeout,
				KeepAlive: 10 * time.Minute,
			}).DialContext(ctx, network, address)
		},
		MaxIdleConns:          int(maxConnections), // 0表示不限
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       0,
		ExpectContinueTimeout: 1 * time.Second,
		TLSHandshakeTimeout:   0, // 不限
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: nil,
	}

	c := &http.Client{
		Timeout:   readTimeout,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	this.locker.Lock()
	this.clientsMap[key] = c
	this.locker.Unlock()

	return c
}

// 重置
func (this *ClientPool) Reset() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.clientsMap = map[string]*http.Client{}
}
