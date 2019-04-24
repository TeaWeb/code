package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
	"strings"
)

type RequestRawRemoteAddrCheckpoint struct {
	Checkpoint
}

func (this *RequestRawRemoteAddrCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	index := strings.LastIndex(req.RemoteAddr, ":")
	if index > -1 {
		value = req.RemoteAddr[:index]
	} else {
		value = req.RemoteAddr
	}
	return
}

func (this *RequestRawRemoteAddrCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
