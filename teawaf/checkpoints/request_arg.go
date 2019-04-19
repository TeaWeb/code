package checkpoints

import "net/http"

type RequestArgCheckPoint struct {
	CheckPoint
}

func (this *RequestArgCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return req.URL.Query().Get(param), nil
}
