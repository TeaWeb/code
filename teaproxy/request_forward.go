package teaproxy

import (
	"bytes"
	"github.com/TeaWeb/code/tealogs/accesslogs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/lists"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

// 正向代理
func (this *Request) Forward(writer *ResponseWriter) error {
	defer this.log()

	// watch
	if this.isWatching {
		// 判断如果Content-Length过长，则截断
		reqData, err := httputil.DumpRequest(this.raw, true)
		if err == nil {
			if len(reqData) > 100240 {
				reqData = reqData[:100240]
			}
			this.requestData = reqData
		}

		writer.SetBodyCopying(true)
	} else {
		max := 512 * 1024 // 512K
		if this.accessLog != nil && lists.ContainsInt(this.accessLog.Fields, accesslogs.AccessLogFieldRequestBody) {
			body, err := ioutil.ReadAll(this.raw.Body)
			if err == nil {
				if len(body) > max {
					this.requestData = body[:max]
				} else {
					this.requestData = body
				}
			}
			this.raw.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		if this.accessLog != nil && lists.ContainsInt(this.accessLog.Fields, accesslogs.AccessLogFieldResponseBody) {
			writer.SetBodyCopying(true)
		}
	}

	if len(this.raw.URL.Scheme) == 0 {
		this.rawScheme = "https"
	}

	this.setProxyHeaders(this.raw.Header)

	if this.method == http.MethodConnect { // connect
		hostConn, err := net.DialTimeout("tcp", this.host, 30*time.Second)
		if err != nil {
			this.serverError(writer)
			this.addError(err)
			return nil
		}

		hijacker, ok := writer.writer.(http.Hijacker)
		if !ok {
			this.serverError(writer)
			this.addError(err)
			return nil
		}

		writer.WriteHeader(http.StatusOK)

		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			this.serverError(writer)
			this.addError(err)
			return nil
		}

		go func() {
			_, _ = io.Copy(clientConn, hostConn)
			_ = clientConn.Close()
			_ = hostConn.Close()
		}()
		go func() {
			_, _ = io.Copy(hostConn, clientConn)
			_ = clientConn.Close()
			_ = hostConn.Close()
		}()
	} else { // http
		this.raw.RequestURI = ""

		// 删除代理相关Header
		for n, _ := range this.raw.Header {
			if lists.ContainsString([]string{"Proxy-Connection", "Connection", "Proxy-Authorization"}, n) {
				this.raw.Header.Del(n)
			}
		}

		client := teautils.SharedHttpClient(30 * time.Second)
		resp, err := client.Do(this.raw)
		if err != nil {
			this.serverError(writer)
			this.addError(err)
			return nil
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		for k, v := range resp.Header {
			if k == "Connection" {
				continue
			}
			for _, subV := range v {
				writer.Header().Add(k, subV)
			}
		}

		writer.Prepare(resp.ContentLength)
		writer.WriteHeader(resp.StatusCode)

		_, _ = io.Copy(writer, resp.Body)
	}

	return nil
}
