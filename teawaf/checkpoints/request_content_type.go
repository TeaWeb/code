package checkpoints

import (
	"net/http"
)

type RequestContentTypeCheckPoint struct {
	CheckPoint
}

func (this *RequestContentTypeCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Header.Get("Content-Type")
	return
}
