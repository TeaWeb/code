package teaproxy

import (
	"errors"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"strings"
	"time"
)

func (this *Request) callURL(writer *ResponseWriter, method string, url string) error {
	req, err := http.NewRequest(method, url, this.raw.Body)
	if err != nil {
		return err
	}

	// 添加当前Header
	req.Header = this.raw.Header

	// 自定义请求Header
	if len(this.requestHeaders) > 0 {
		for _, header := range this.requestHeaders {
			if !header.On {
				continue
			}
			if header.HasVariables() {
				req.Header.Set(header.Name, this.Format(header.Value))
			} else {
				req.Header.Set(header.Name, header.Value)
			}
		}
	}

	// 代理头部
	this.setProxyHeaders(req.Header)

	var client *http.Client = nil
	if len(req.Host) > 0 {
		host := req.Host
		if !strings.Contains(host, ":") {
			if req.URL.Scheme == "https" {
				host += ":443"
			} else {
				host += ":80"
			}
		}
		client = SharedClientPool.client("", host, 60*time.Second, 0, 0)
	} else {
		client = teautils.SharedHttpClient(60 * time.Second)
	}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error(errors.New(req.URL.String() + ": " + err.Error()))
		this.addError(err)
		this.serverError(writer)
		return err
	}
	defer resp.Body.Close()

	// Header
	this.WriteResponseHeaders(writer, resp.StatusCode)

	writer.AddHeaders(resp.Header)
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	// 输出内容
	buf := make([]byte, 1024) // TODO 可以配置
	_, err = io.CopyBuffer(writer, resp.Body, buf)

	return err
}
