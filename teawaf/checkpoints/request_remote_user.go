package checkpoints

import (
	"net/http"
)

type RequestRemoteUserCheckPoint struct {
	CheckPoint
}

func (this *RequestRemoteUserCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	username, _, ok := req.BasicAuth()
	if !ok {
		value = ""
		return
	}
	value = username
	return
}
