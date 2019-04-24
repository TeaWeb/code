package teaproxy

import (
	"errors"
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

	// ip
	remoteAddr := this.requestRemoteAddr()
	if len(remoteAddr) > 0 {
		index := strings.Index(this.raw.RemoteAddr, ":")
		ip := ""
		if index > -1 {
			ip = this.raw.RemoteAddr[:index]
		} else {
			ip = this.raw.RemoteAddr
		}
		req.Header["X-Real-IP"] = []string{ip}
		req.Header.Set("X-Forwarded-For", ip)
		req.Header.Set("X-Forwarded-By", ip)
	}

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
		client = SharedClientPool.client("", host, 30*time.Second, 0, 0)
	} else {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
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
	writer.AddHeaders(resp.Header)
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	// 输出内容
	_, err = io.Copy(writer, resp.Body)

	return err
}
