package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type RequestCookieCheckpoint struct {
	Checkpoint
}

func (this *RequestCookieCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	cookie, err := req.Cookie(param)
	if err != nil {
		value = ""
		return
	}

	value = cookie.Value
	return
}

func (this *RequestCookieCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
