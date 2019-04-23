package checkpoints

import (
	"net/http"
)

type RequestSchemeCheckpoint struct {
	Checkpoint
}

func (this *RequestSchemeCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.URL.Scheme
	return
}
