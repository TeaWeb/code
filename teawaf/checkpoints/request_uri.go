package checkpoints

import (
	"net/http"
)

type RequestURICheckpoint struct {
	Checkpoint
}

func (this *RequestURICheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.RequestURI
	return
}
