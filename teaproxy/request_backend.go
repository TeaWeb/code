package teaproxy

import (
	"context"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/url"
	"strings"
	"time"
)

// 调用后端服务器
func (this *Request) callBackend(writer *ResponseWriter) error {
	// 是否为websocket请求
	if this.raw.Header.Get("Upgrade") == "websocket" {
		websocket := teaconfigs.NewWebsocketConfig()
		websocket.On = true
		websocket.AllowAllOrigins = true
		websocket.ForwardMode = teaconfigs.WebsocketForwardModeWebsocket
		websocket.Backends = []*teaconfigs.BackendConfig{this.backend}
		websocket.Origins = []string{}
		websocket.HandshakeTimeout = "5s"
		this.websocket = websocket
		err := websocket.Validate()
		if err != nil {
			logs.Error(err)
		}
		return this.callWebsocket(writer)
	}

	this.backend.IncreaseConn()
	defer this.backend.DecreaseConn()

	if len(this.backend.Address) == 0 {
		this.serverError(writer)
		logs.Error(errors.New("backend address should not be empty"))
		this.addError(errors.New("backend address should not be empty"))
		return nil
	}

	this.raw.URL.Host = this.host

	if this.backend.HasHost() {
		this.raw.Host = this.Format(this.backend.Host)
	}

	if len(this.backend.Scheme) > 0 && this.backend.Scheme != "http" {
		this.raw.URL.Scheme = this.backend.Scheme
	} else {
		this.raw.URL.Scheme = this.scheme
	}

	// new uri
	if this.backend.HasRequestURI() {
		uri := this.Format(this.backend.RequestPath())

		u, err := url.ParseRequestURI(uri)
		if err == nil {
			this.raw.URL.Path = u.Path
			this.raw.URL.RawQuery = u.RawQuery

			args := this.Format(this.backend.RequestArgs())
			if len(args) > 0 {
				if len(u.RawQuery) > 0 {
					this.raw.URL.RawQuery += "&" + args
				} else {
					this.raw.URL.RawQuery += args
				}
			}
		}
	} else {
		u, err := url.ParseRequestURI(this.uri)
		if err == nil {
			this.raw.URL.Path = u.Path
			this.raw.URL.RawQuery = u.RawQuery
		}
	}

	// 设置代理相关的头部
	// 参考 https://tools.ietf.org/html/rfc7239
	this.setProxyHeaders(this.raw.Header)

	this.raw.Header.Set("Connection", "keep-alive")

	// 自定义请求Header
	if len(this.requestHeaders) > 0 {
		for _, header := range this.requestHeaders {
			if !header.On {
				continue
			}
			if header.HasVariables() {
				this.raw.Header.Set(header.Name, this.Format(header.Value))
			} else {
				this.raw.Header.Set(header.Name, header.Value)
			}

			// 支持修改Host
			if header.Name == "Host" {
				this.raw.Host = header.Value
			}
		}
	}

	client := SharedClientPool.client(this.backend.Id, this.backend.Address, this.backend.FailTimeoutDuration(), this.backend.ReadTimeoutDuration(), this.backend.MaxConns)

	this.raw.RequestURI = ""

	resp, err := client.Do(this.raw)
	if err != nil {
		// 客户端取消请求，则不提示
		httpErr, ok := err.(*url.Error)
		if !ok || httpErr.Err != context.Canceled {
			// 如果超过最大失败次数，则下线
			if !this.backend.HasCheckURL() {
				currentFails := this.backend.IncreaseFails()
				if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
					this.backend.IsDown = true
					this.backend.DownTime = time.Now()
					if this.websocket != nil {
						this.websocket.SetupScheduling(false)
					} else {
						this.server.SetupScheduling(false)
					}
				}
			}

			this.serverError(writer)

			logs.Println("[proxy]" + err.Error())
			this.addError(err)
		} else {
			this.serverError(writer)
			this.addError(err)
		}
		return nil
	}

	// waf
	if this.waf != nil {
		if this.callWAFResponse(resp, writer) {
			resp.Body.Close()
			return nil
		}
	}

	data := []byte{}
	bodyRead := false
	if resp.ContentLength > 0 && resp.ContentLength < 2048 { // 内容比较少的直接读取，以加快响应速度
		bodyRead = true

		buf := make([]byte, 256)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				data = append(data, buf[:n]...)
			}
			if err != nil {
				break
			}
		}

		resp.ContentLength = int64(len(data))
		resp.Body.Close()
	} else {
		defer resp.Body.Close()
	}

	// 清除错误次数
	if resp.StatusCode >= 200 && !this.backend.HasCheckURL() {
		if !this.backend.IsDown && this.backend.CurrentFails > 0 {
			this.backend.CurrentFails = 0
		}
	}

	// 特殊页面
	if len(this.pages) > 0 && this.callPage(writer, resp.StatusCode) {
		return nil
	}

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 设置Header
	hasCharset := len(this.charset) > 0
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(k)) {
			continue
		}
		for _, subV := range v {
			// 字符集
			if hasCharset && k == "Content-Type" {
				if _, found := textMimeMap[subV]; found {
					if !strings.Contains(subV, "charset=") {
						subV += "; charset=" + this.charset
					}
				}
			}

			writer.Header().Add(k, subV)
		}
	}

	// 自定义响应Headers
	this.WriteResponseHeaders(writer, resp.StatusCode)

	// 响应回调
	if this.responseCallback != nil {
		this.responseCallback(writer)
	}

	// 准备
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	if bodyRead {
		_, err = writer.Write(data)
	} else {
		_, err = io.Copy(writer, resp.Body)
	}
	if err != nil {
		logs.Error(err)
		this.addError(err)
		return nil
	}
	return nil
}
