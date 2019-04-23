package checkpoints

import (
	"net/http"
)

type RequestLengthCheckpoint struct {
	Checkpoint
}

func (this *RequestLengthCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.ContentLength
	return
}
