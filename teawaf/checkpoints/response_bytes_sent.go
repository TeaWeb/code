package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

// ${bytesSent}
type ResponseBytesSentCheckpoint struct {
	Checkpoint
}

func (this *ResponseBytesSentCheckpoint) IsRequest() bool {
	return false
}

func (this *ResponseBytesSentCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = 0
	return
}

func (this *ResponseBytesSentCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	value = 0
	if resp != nil {
		value = resp.ContentLength
	}
	return
}
