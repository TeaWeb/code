package checkpoints

import (
	"net/http"
)

type RequestContentTypeCheckpoint struct {
	Checkpoint
}

func (this *RequestContentTypeCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Header.Get("Content-Type")
	return
}
