package checkpoints

import (
	"net/http"
)

type RequestArgsCheckPoint struct {
	CheckPoint
}

func (this *RequestArgsCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.URL.RawQuery
	return
}
