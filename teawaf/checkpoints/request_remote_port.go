package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"strings"
)

type RequestRemotePortCheckpoint struct {
	Checkpoint
}

func (this *RequestRemotePortCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	remoteAddr := req.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index < 0 {
		value = 0
	} else {
		value = remoteAddr[index+1:]
	}
	return
}

func (this *RequestRemotePortCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
