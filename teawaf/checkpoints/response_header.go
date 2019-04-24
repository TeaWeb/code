package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

// ${responseHeader.arg}
type ResponseHeaderCheckpoint struct {
	Checkpoint
}

func (this *ResponseHeaderCheckpoint) IsRequest() bool {
	return false
}

func (this *ResponseHeaderCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = ""
	return
}

func (this *ResponseHeaderCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if resp != nil && resp.Header != nil {
		value = resp.Header.Get(param)
	} else {
		value = ""
	}
	return
}
