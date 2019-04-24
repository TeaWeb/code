package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

// just a sample checkpoint, copy and change it for your new checkpoint
type SampleRequestCheckpoint struct {
	Checkpoint
}

func (this *SampleRequestCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	return
}

func (this *SampleRequestCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
