package teautils

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func DumpResponse(resp *http.Response) (header []byte, body []byte, err error) {
	header, err = httputil.DumpResponse(resp, false)
	body, err = ioutil.ReadAll(resp.Body)
	return
}
