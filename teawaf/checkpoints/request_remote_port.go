package checkpoints

import (
	"net/http"
	"strings"
)

type RequestRemotePortCheckPoint struct {
	CheckPoint
}

func (this *RequestRemotePortCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	remoteAddr := req.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index < 0 {
		value = 0
	} else {
		value = remoteAddr[index+1:]
	}
	return
}
