package teaproxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/gorilla/websocket"
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 调用Websocket
func (this *Request) callWebsocket(writer *ResponseWriter) error {
	if this.backend == nil {
		err := errors.New(this.requestPath() + ": no available backends for websocket")
		logs.Error(err)
		this.addError(err)
		this.serverError(writer)
		return err
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: this.websocket.HandshakeTimeoutDuration(),
		CheckOrigin: func(r *http.Request) bool {
			if this.websocket.AllowAllOrigins {
				return true
			}
			origin := r.Header.Get("Origin")
			if len(origin) == 0 {
				return false
			}
			return this.websocket.MatchOrigin(origin)
		},
	}

	// 自动补充Header
	this.raw.Header.Set("Connection", "upgrade")
	if len(this.raw.Header.Get("Upgrade")) == 0 {
		this.raw.Header.Set("Upgrade", "websocket")
	}

	// 接收客户端连接
	client, err := upgrader.Upgrade(this.responseWriter.Raw(), this.raw, nil)
	if err != nil {
		logs.Error(errors.New("upgrade: " + err.Error()))
		this.addError(errors.New("upgrade: " + err.Error()))
		return err
	}
	defer client.Close()

	if this.websocket.ForwardMode == teaconfigs.WebsocketForwardModeWebsocket {
		// 判断最大连接数
		if this.backend.MaxConns > 0 && this.backend.CurrentConns >= this.backend.MaxConns {
			this.serverError(writer)
			logs.Error(errors.New("too many connections"))
			this.addError(errors.New("too many connections"))
			return nil
		}

		// 增加连接数
		this.backend.IncreaseConn()
		defer this.backend.DecreaseConn()

		// 连接后端服务器
		scheme := "ws"
		if this.backend.Scheme == "https" {
			scheme = "wss"
		}
		wsURL := url.URL{Scheme: scheme, Host: this.backend.Address, Path: this.raw.RequestURI}
		dialer := websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: this.backend.FailTimeoutDuration(),
		}
		header := http.Header{}
		{
			origin, ok := this.raw.Header["Origin"]
			if ok {
				header["Origin"] = origin
			}
		}

		// 自定义请求Header
		for _, h := range this.requestHeaders {
			if !h.On {
				continue
			}
			if h.HasVariables() {
				header[h.Name] = []string{this.Format(h.Value)}
			} else {
				header[h.Name] = []string{h.Value}
			}
		}

		server, _, err := dialer.Dial(wsURL.String(), header)
		if err != nil {
			logs.Error(err)
			this.addError(err)
			currentFails := this.backend.IncreaseFails()
			if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
				this.backend.IsDown = true
				this.backend.DownTime = time.Now()
				this.websocket.SetupScheduling(false)
			}
			return err
		}
		defer server.Close()

		// 设置关闭连接的处理函数
		clientIsClosed := false
		serverIsClosed := false
		client.SetCloseHandler(func(code int, text string) error {
			if serverIsClosed {
				return nil
			}
			serverIsClosed = true
			return server.Close()
		})

		// 从客户端接收数据
		go func() {
			for {
				messageType, message, err := client.ReadMessage()
				if err != nil {
					closeErr, ok := err.(*websocket.CloseError)
					if !ok && closeErr != nil && closeErr.Code != websocket.CloseGoingAway {
						logs.Error(err)
						this.addError(err)
					}
					clientIsClosed = true
					break
				}
				server.WriteMessage(messageType, message)
			}
		}()

		// 从后端服务器读取数据
		for {
			messageType, message, err := server.ReadMessage()
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
				if !ok && closeErr != nil && closeErr.Code != websocket.CloseGoingAway {
					logs.Error(err)
					this.addError(err)
				}
				serverIsClosed = true
				server.Close()
				if !clientIsClosed {
					client.Close()
				}
				break
			}
			client.WriteMessage(messageType, message)
		}
	} else if this.websocket.ForwardMode == teaconfigs.WebsocketForwardModeHttp {
		messageQueue := make(chan []byte, 1024)
		quit := make(chan bool)
		go func() {
		FOR:
			for {
				select {
				case message := <-messageQueue:
					{
						this.raw.Method = http.MethodPut
						responseWriter := NewResponseWriter(nil)
						responseWriter.SetBodyCopying(true)
						this.raw.Body = ioutil.NopCloser(bytes.NewReader(message))
						this.raw.Header.Del("Upgrade")
						err := this.callBackend(responseWriter)
						if err != nil {
							continue FOR
						}
						if responseWriter.StatusCode() != http.StatusOK {
							logs.Error(errors.New(this.requestURI() + ": invalid response from backend: " + fmt.Sprintf("%d", responseWriter.StatusCode()) + " " + http.StatusText(responseWriter.StatusCode())))
							this.addError(errors.New(this.requestURI() + ": invalid response from backend: " + fmt.Sprintf("%d", responseWriter.StatusCode()) + " " + http.StatusText(responseWriter.StatusCode())))
							continue FOR
						}
						client.WriteMessage(websocket.TextMessage, responseWriter.Body())
					}
				case <-quit:
					break FOR
				}
			}
		}()
		for {
			messageType, message, err := client.ReadMessage()
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
				if !ok || closeErr.Code != websocket.CloseGoingAway {
					logs.Error(err)
					this.addError(err)
				}
				quit <- true
				break
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				messageQueue <- message
			}
		}
	}

	return nil
}
