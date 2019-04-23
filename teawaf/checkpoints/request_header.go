package checkpoints

import (
	"net/http"
	"strings"
)

type RequestHeaderCheckpoint struct {
	Checkpoint
}

func (this *RequestHeaderCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	v, found := req.Header[param]
	if !found {
		value = ""
		return
	}
	value = strings.Join(v, ";")
	return
}
