package teaproxy

import (
	"bufio"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
)

// 隧道连接
type TunnelConnection struct {
	conn   net.Conn
	reader *bufio.Reader
	locker *sync.Mutex
}

// 获取新对象
func NewTunnelConnection(conn net.Conn) *TunnelConnection {
	return &TunnelConnection{
		conn:   conn,
		reader: bufio.NewReader(conn),
		locker: &sync.Mutex{},
	}
}

// 发送请求
func (this *TunnelConnection) Write(req *http.Request) (*http.Response, error) {
	if this.reader == nil {
		return nil, errors.New("[tunnel]no tunnel reader")
	}

	this.locker.Lock()

	data, err := httputil.DumpRequest(req, true)
	_, err = this.conn.Write(data)
	if err != nil {
		this.locker.Unlock()
		return nil, err
	}

	resp, err := http.ReadResponse(this.reader, req)
	if err != nil {
		this.locker.Unlock()
		return resp, err
	}
	resp.Body = &TunnelResponseBody{
		ReadCloser: resp.Body,
		locker:     this.locker,
	}
	return resp, nil
}

// 关闭
func (this *TunnelConnection) Close() error {
	return this.conn.Close()
}
