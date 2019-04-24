package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type RequestLengthCheckpoint struct {
	Checkpoint
}

func (this *RequestLengthCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = req.ContentLength
	return
}

func (this *RequestLengthCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
