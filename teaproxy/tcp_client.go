package teaproxy

import (
	"crypto/tls"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net"
	"sync/atomic"
	"time"
)

const (
	TCPClientMaxAttempts = 3    // 失败最多尝试次数
	TCPClientBufferSize  = 1024 // 缓冲区尺寸（字节）
	TCPClientStreamSize  = 128  // 读取客户端数据的队列长度
)

// TCP连接客户端
type TCPClient struct {
	serverPool     func() *teaconfigs.ServerConfig
	lConn          net.Conn
	stream         chan []byte
	streamIsClosed bool
	initData       []byte

	lActive bool

	excludingBackendIds []string
	backend             *teaconfigs.BackendConfig
	rConn               net.Conn

	readSpeed  int64
	writeSpeed int64
}

// 创建新的客户端对象
func NewTCPClient(serverPool func() *teaconfigs.ServerConfig, conn net.Conn) *TCPClient {
	return &TCPClient{
		serverPool: serverPool,
		lConn:      conn,
		stream:     make(chan []byte, TCPClientStreamSize),
		lActive:    true,
	}
}

// 获取左连接 - 客户端
func (this *TCPClient) LConn() net.Conn {
	return this.lConn
}

// 获取右连接 - 后端服务器
func (this *TCPClient) RConn() net.Conn {
	return this.rConn
}

// 连接后端服务器
func (this *TCPClient) Connect() {
	ticker := teautils.Every(1*time.Second, func(ticker *teautils.Ticker) {
		atomic.StoreInt64(&this.readSpeed, 0)
		atomic.StoreInt64(&this.writeSpeed, 0)
	})

	go this.connect()
	this.read()

	ticker.Stop()
}

// 关闭
func (this *TCPClient) Close() error {
	this.lActive = false
	lCloseError := this.lConn.Close()

	if this.rConn != nil {
		this.rConn.Close()
	}

	// 关闭stream
	if !this.streamIsClosed {
		this.streamIsClosed = true
		close(this.stream)
	}

	return lCloseError
}

// 获取读取的速度
func (this *TCPClient) ReadSpeed() int64 {
	return atomic.LoadInt64(&this.readSpeed)
}

// 获取写入的速度
func (this *TCPClient) WriteSpeed() int64 {
	return atomic.LoadInt64(&this.writeSpeed)
}

// 连接后端服务器
func (this *TCPClient) connect() {
	if this.serverPool == nil {
		logs.Error(errors.New("'serverPool' must not be nil"))
		return
	}

	server := this.serverPool()
	if server == nil {
		logs.Error(errors.New("no available server for the connection"))
		return
	}

	if !this.lActive {
		return
	}

	requestCall := shared.NewRequestCall()

	for i := 0; i < TCPClientMaxAttempts; i++ {
		// 查找下一个Backend
		if len(this.excludingBackendIds) == 0 {
			this.backend = server.NextBackend(requestCall)
		} else {
			this.backend = server.NextBackendIgnore(requestCall, this.excludingBackendIds)
		}
		if this.backend == nil {
			logs.Println("[proxy][tcp]no available backends for server '" + server.Description)
			this.Close()
			break
		}

		// 是否超过最大连接数
		currentConns := this.backend.IncreaseConn()
		if this.backend.MaxConns > 0 && currentConns > this.backend.MaxConns {
			this.fail(server, errors.New("too many connections"))
			this.backend.DecreaseConn()
			continue
		}

		// 开始连接
		switch this.backend.Scheme {
		case "tcp":
			conn, err := net.DialTimeout("tcp", this.backend.Address, this.backend.FailTimeoutDuration())
			if err != nil {
				this.error(server, err)
				this.backend.DecreaseConn()
				break
			}
			this.rConn = conn
		case "tcp+tls":
			conn, err := tls.DialWithDialer(&net.Dialer{
				Timeout: this.backend.FailTimeoutDuration(),
			}, "tcp", this.backend.Address, &tls.Config{
				InsecureSkipVerify: true,
			})
			if err != nil {
				this.error(server, err)
				this.backend.DecreaseConn()
				break
			}
			this.rConn = conn
		}

		// 已连接就不继续尝试
		if this.rConn != nil {
			break
		}
	}

	// 没连接到后端就中断
	if this.rConn == nil {
		this.Close()
		return
	}

	// 写入
	go func() {
		if len(this.initData) > 0 {
			_, err := this.rConn.Write(this.initData)
			if err == nil {
				this.initData = nil
			}
		}

		for data := range this.stream {
			if data == nil {
				break
			}
			_, err := this.rConn.Write(data)
			if err != nil {
				if server.TCP.FailResend {
					this.initData = data
				}
				break
			}
		}
	}()

	// 读取
	buf := make([]byte, TCPClientBufferSize)
	for {
		n, err := this.rConn.Read(buf)
		if n > 0 {
			if this.lActive {
				this.lConn.Write(buf[:n])
				atomic.AddInt64(&this.readSpeed, int64(n))
			}
		}
		if err != nil {
			if this.lActive {
				this.error(server, err)
			}
			break
		}
	}

	this.backend.DecreaseConn()

	// 重新连接
	if this.lActive && (len(this.excludingBackendIds) == 0 || server.TCP.FailReconnect) {
		logs.Println("[proxy][tcp]reconnecting another backend for client '" + this.lConn.RemoteAddr().String() + "'")
		this.backend = nil
		this.rConn = nil
		if !this.streamIsClosed {
			close(this.stream)
			this.stream = make(chan []byte, TCPClientStreamSize)
		}
		this.connect()
	}
}

// 读取客户端数据
func (this *TCPClient) read() {
	buf := make([]byte, TCPClientBufferSize)
	for {
		n, err := this.lConn.Read(buf)
		if n > 0 {
			this.stream <- buf[:n]

			atomic.AddInt64(&this.writeSpeed, int64(n))
		}
		if err != nil {
			this.lActive = false
			this.Close()
			break
		}
	}
}

// 处理连接失败
func (this *TCPClient) fail(server *teaconfigs.ServerConfig, err error) {
	if this.backend == nil {
		return
	}

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		logs.Println("[proxy][tcp]lost connection to backend '" + this.backend.Address + "'" + " for server '" + server.Description + "'")
	} else {
		logs.Println("[proxy][tcp]failed to connect backend '" + this.backend.Address + "'" + " for server '" + server.Description + "': " + err.Error())
	}

	this.excludingBackendIds = append(this.excludingBackendIds, this.backend.Id)
}

// 处理连接错误
func (this *TCPClient) error(server *teaconfigs.ServerConfig, err error) {
	if this.backend == nil {
		return
	}

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		logs.Println("[proxy][tcp]lost connection to backend '" + this.backend.Address + "'" + " for server '" + server.Description + "'")
	} else {
		logs.Println("[proxy][tcp]failed to connect backend '" + this.backend.Address + "'" + " for server '" + server.Description + "': " + err.Error())
	}

	this.excludingBackendIds = append(this.excludingBackendIds, this.backend.Id)

	currentFails := this.backend.IncreaseFails()
	if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
		this.backend.IsDown = true
		this.backend.DownTime = time.Now()

		// 下线通知
		noticeutils.NotifyProxyBackendDownMessage(server.Id, this.backend, nil, nil)

		server.SetupScheduling(false)
	}
}
