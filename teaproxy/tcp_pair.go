package teaproxy

import (
	"net"
)

// 创建一对TCP连接，可以相互拷贝内容
type TCPPair struct {
	lConn net.Conn
	rConn net.Conn
}

// 创建新的TCP连接对
func NewTCPPair(lConn, rConn net.Conn) *TCPPair {
	return &TCPPair{
		lConn: lConn,
		rConn: rConn,
	}
}

// 开始传送
func (this *TCPPair) Transfer() error {
	// l -> r
	go func() {
		buf := make([]byte, 256)
		for {
			n, err := this.lConn.Read(buf)
			if n > 0 {
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
			_, err = this.lConn.Write(buf[:n])
			if err != nil {
				this.Close()
				break
			}
		}
		if err != nil {
			this.Close()
			break
		}
	}

	return nil
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
