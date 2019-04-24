package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

// ${bytesSent}
type ResponseStatusCheckpoint struct {
	Checkpoint
}

func (this *ResponseStatusCheckpoint) IsRequest() bool {
	return false
}

func (this *ResponseStatusCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = 0
	return
}

func (this *ResponseStatusCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if resp != nil {
		value = resp.StatusCode
	}
	return
}
