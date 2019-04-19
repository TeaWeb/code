package checkpoints

import "net/http"

type CheckPoint struct {
}

func (this *CheckPoint) IsRequest() bool {
	return true
}

func (this *CheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return
}

func (this *CheckPoint) ResponseValue(req *http.Request, resp *http.Response, param string) (value interface{}, err error) {
	return
}
