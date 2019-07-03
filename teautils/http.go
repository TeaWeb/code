package teautils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

// 导出响应
func DumpResponse(resp *http.Response) (header []byte, body []byte, err error) {
	header, err = httputil.DumpResponse(resp, false)
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// 获取一个Client
func NewHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:          1024,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       0,
			ExpectContinueTimeout: 1 * time.Second,
			TLSHandshakeTimeout:   0,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
