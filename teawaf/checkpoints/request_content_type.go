package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type RequestContentTypeCheckpoint struct {
	Checkpoint
}

func (this *RequestContentTypeCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = req.Header.Get("Content-Type")
	return
}

func (this *RequestContentTypeCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
