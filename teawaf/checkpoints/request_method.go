package checkpoints

import (
	"net/http"
)

type RequestMethodCheckpoint struct {
	Checkpoint
}

func (this *RequestMethodCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Method
	return
}
