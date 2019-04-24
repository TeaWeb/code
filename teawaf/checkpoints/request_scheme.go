package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type RequestSchemeCheckpoint struct {
	Checkpoint
}

func (this *RequestSchemeCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = req.URL.Scheme
	return
}

func (this *RequestSchemeCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
