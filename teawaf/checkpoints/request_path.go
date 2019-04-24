package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type RequestPathCheckpoint struct {
	Checkpoint
}

func (this *RequestPathCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	return req.URL.Path, nil, nil
}

func (this *RequestPathCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
