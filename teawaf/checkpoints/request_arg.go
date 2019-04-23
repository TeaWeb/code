package checkpoints

import "net/http"

type RequestArgCheckpoint struct {
	Checkpoint
}

func (this *RequestArgCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return req.URL.Query().Get(param), nil
}
