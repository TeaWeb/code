package checkpoints

import (
	"net/http"
)

type RequestHostCheckpoint struct {
	Checkpoint
}

func (this *RequestHostCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Host
	return
}
