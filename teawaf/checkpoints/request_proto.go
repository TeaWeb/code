package checkpoints

import (
	"net/http"
)

type RequestProtoCheckpoint struct {
	Checkpoint
}

func (this *RequestProtoCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Proto
	return
}
