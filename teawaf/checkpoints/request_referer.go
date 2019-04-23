package checkpoints

import (
	"net/http"
)

type RequestRefererCheckpoint struct {
	Checkpoint
}

func (this *RequestRefererCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Referer()
	return
}
