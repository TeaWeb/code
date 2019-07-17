package teaproxy

import (
	"crypto/tls"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/logs"
	"net"
	"time"
)

// TCP后端连接管理
type TCPBackend struct {
	server      *teaconfigs.ServerConfig
	backend     *teaconfigs.BackendConfig
	requestCall *shared.RequestCall
	pair        *TCPPair

	ignoreBackendIds []string

	initData    []byte
	clientConn  net.Conn
	backendConn net.Conn

	isInterrupted bool
	attempts      int

	successHandler func(pair *TCPPair)
	failHandler    func(pair *TCPPair)
}

// 创建新对象
func NewTCPBackend(server *teaconfigs.ServerConfig, clientConn net.Conn) *TCPBackend {
	return &TCPBackend{
		server:      server,
		requestCall: shared.NewRequestCall(),
		clientConn:  clientConn,
	}
}

// 成功时回调
func (this *TCPBackend) OnSuccess(successHandler func(pair *TCPPair)) {
	this.successHandler = successHandler
}

// 失败时回调
func (this *TCPBackend) OnFail(failHandler func(pair *TCPPair)) {
	this.failHandler = failHandler
}

// 连接
func (this *TCPBackend) Connect() {
	this.backend = this.server.NextBackend(this.requestCall)
	if this.backend == nil {
		this.clientConn.Close()
		logs.Println("[proxy][tcp]no backends for '" + this.server.Description + "'")
		return
	}
	this.connect()
}

// 实际的连接动作
func (this *TCPBackend) connect() (err error) {
	currentConns := this.backend.IncreaseConn()
	defer this.backend.DecreaseConn()

	// 是否超过最大连接数
	if this.backend.MaxConns > 0 && currentConns > this.backend.MaxConns {
		this.fail(errors.New("too many connections"))
		return
	}

	if this.backend.Scheme == "tcp" || len(this.backend.Scheme) == 0 { // TCP
		this.backendConn, err = net.DialTimeout("tcp", this.backend.Address, this.backend.FailTimeoutDuration())
		if err != nil {
			this.fail(err)
			return
		}
	} else if this.backend.Scheme == "tcp+tls" { // TCP+TLS
		this.backendConn, err = tls.Dial("tcp", this.backend.Address, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			this.fail(err)
			return
		}
	} else { // neither tcp nor tcp+tls
		this.fail(errors.New("invalid scheme"))
		return
	}

	// 如果能够连接，则重置统计数字
	this.attempts = 0

	// 写入初始化数据
	if len(this.initData) > 0 {
		_, err = this.backendConn.Write(this.initData)
		if err != nil {
			this.fail(err)
			return
		}
	}

	// 创建传输对
	if this.pair == nil {
		this.pair = NewTCPPair(this.clientConn, this.backendConn)
	} else {
		this.pair.SetRConn(this.backendConn)
	}
	if this.server.TCP.FailReconnect {
		this.pair.OnRightDisconnect(func(data []byte) {
			if this.server.TCP.FailResend {
				this.initData = data
			}
			this.isInterrupted = true
			this.fail(errors.New("connection interrupted"))
		})
	}

	// 成功回调
	if this.successHandler != nil {
		this.successHandler(this.pair)
	}

	// 开始传输并阻塞当前程序
	err = this.pair.Transfer()
	if err != nil {
		logs.Println("[proxy][tcp]" + err.Error())

		// 失败回调
		if this.failHandler != nil {
			this.failHandler(this.pair)
		}
		return
	}

	// 失败回调
	if this.failHandler != nil {
		this.failHandler(this.pair)
	}

	return
}

// 处理连接失败
func (this *TCPBackend) fail(err error) {
	if this.backend != nil {
		this.ignoreBackendIds = append(this.ignoreBackendIds, this.backend.Id)
	}

	// 最多尝试连接到后端服务器次数
	const MaxAttempts = 3

	// 超出最大尝试次数，则关闭客户端
	if this.attempts > MaxAttempts {
		this.clientConn.Close()
		return
	}

	logs.Println("[proxy][tcp]failed to connect backend '" + this.backend.Address + "'" + " for server '" + this.server.Description + ": " + err.Error())

	currentFails := this.backend.IncreaseFails()
	if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
		this.backend.IsDown = true
		this.backend.DownTime = time.Now()
		this.server.SetupScheduling(false)
	}

	// 查找下一个
	if this.isInterrupted && !this.server.TCP.FailReconnect {
		return
	}
	this.backend = this.server.NextBackendIgnore(this.requestCall, this.ignoreBackendIds)
	if this.backend == nil {
		this.clientConn.Close()
		return
	}

	this.attempts++
	go this.connect()
}
