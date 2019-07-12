package teaproxy

import (
	"github.com/iwind/TeaGo/timers"
	"net"
	"sync/atomic"
	"time"
)

// 创建一对TCP连接，可以相互拷贝内容
type TCPPair struct {
	lConn net.Conn
	rConn net.Conn

	rBytesSpeed1s int64 // 1秒钟之内读取的字节数
	wBytesSpeed1s int64 // 1秒钟之内写入的字节数
}

// 创建新的TCP连接对
func NewTCPPair(lConn, rConn net.Conn) *TCPPair {
	return &TCPPair{
		lConn: lConn,
		rConn: rConn,
	}
}

// 左连接
func (this *TCPPair) LConn() net.Conn {
	return this.lConn
}

// 右连接
func (this *TCPPair) RConn() net.Conn {
	return this.rConn
}

// 开始传送
func (this *TCPPair) Transfer() error {
	// 每一秒钟清除一下速度数据
	timers.Every(1*time.Second, func(ticker *time.Ticker) {
		atomic.StoreInt64(&this.rBytesSpeed1s, 0)
		atomic.StoreInt64(&this.wBytesSpeed1s, 0)
	})

	// l -> r
	go func() {
		buf := make([]byte, 256)
		for {
			n, err := this.lConn.Read(buf)
			if n > 0 {
				atomic.AddInt64(&this.wBytesSpeed1s, int64(n))

				_, err = this.rConn.Write(buf[:n])
				if err != nil {
					this.Close()
					return
				}
			}
			if err != nil {
				this.Close()
				return
			}
		}
	}()

	// l <- r
	// 此时不用go routine，是为了hold住协程
	buf := make([]byte, 256)
	for {
		n, err := this.rConn.Read(buf)
		if n > 0 {
			n2, err := this.lConn.Write(buf[:n])
			if err != nil {
				this.Close()
				break
			}

			atomic.AddInt64(&this.rBytesSpeed1s, int64(n2))
		}
		if err != nil {
			this.Close()
			break
		}
	}

	return nil
}

// 客户端读取速度
func (this *TCPPair) ReadSpeed() int64 {
	return atomic.LoadInt64(&this.rBytesSpeed1s)
}

// 客户端写入速度
func (this *TCPPair) WriteSpeed() int64 {
	return atomic.LoadInt64(&this.wBytesSpeed1s)
}

// 关闭
func (this *TCPPair) Close() error {
	err1 := this.lConn.Close()
	err2 := this.rConn.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
