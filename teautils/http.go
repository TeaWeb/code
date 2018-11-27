package teautils

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

// 导出响应
func DumpResponse(resp *http.Response) (header []byte, body []byte, err error) {
	header, err = httputil.DumpResponse(resp, false)
	body, err = ioutil.ReadAll(resp.Body)
	return
}
