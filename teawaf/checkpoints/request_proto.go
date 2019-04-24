package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type RequestProtoCheckpoint struct {
	Checkpoint
}

func (this *RequestProtoCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = req.Proto
	return
}

func (this *RequestProtoCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
