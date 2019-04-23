package checkpoints

import (
	"net/http"
)

type RequestArgsCheckpoint struct {
	Checkpoint
}

func (this *RequestArgsCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.URL.RawQuery
	return
}
