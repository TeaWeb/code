package checkpoints

import "net/http"

type RequestPathCheckpoint struct {
	Checkpoint
}

func (this *RequestPathCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return req.URL.Path, nil
}
