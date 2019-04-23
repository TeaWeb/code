package checkpoints

import "net/http"

type Checkpoint struct {
}

func (this *Checkpoint) IsRequest() bool {
	return true
}

func (this *Checkpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return
}

func (this *Checkpoint) ResponseValue(req *http.Request, resp *http.Response, param string) (value interface{}, err error) {
	return
}
