package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
)

type RequestArgCheckpoint struct {
	Checkpoint
}

func (this *RequestArgCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	return req.URL.Query().Get(param), nil, nil
}

func (this *RequestArgCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
