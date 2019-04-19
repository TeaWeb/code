package checkpoints

import (
	"net/http"
	"strings"
)

type RequestRawRemoteAddrCheckPoint struct {
	CheckPoint
}

func (this *RequestRawRemoteAddrCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	index := strings.LastIndex(req.RemoteAddr, ":")
	if index > -1 {
		value = req.RemoteAddr[:index]
	} else {
		value = req.RemoteAddr
	}
	return
}
