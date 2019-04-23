package checkpoints

import (
	"net/http"
)

type RequestRemoteUserCheckpoint struct {
	Checkpoint
}

func (this *RequestRemoteUserCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	username, _, ok := req.BasicAuth()
	if !ok {
		value = ""
		return
	}
	value = username
	return
}
