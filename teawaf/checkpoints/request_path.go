package checkpoints

import "net/http"

type RequestPathCheckPoint struct {
	CheckPoint
}

func (this *RequestPathCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return req.URL.Path, nil
}
