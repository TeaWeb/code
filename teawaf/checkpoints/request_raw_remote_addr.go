package checkpoints

import (
	"net/http"
	"strings"
)

type RequestRawRemoteAddrCheckpoint struct {
	Checkpoint
}

func (this *RequestRawRemoteAddrCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	index := strings.LastIndex(req.RemoteAddr, ":")
	if index > -1 {
		value = req.RemoteAddr[:index]
	} else {
		value = req.RemoteAddr
	}
	return
}
