package checkpoints

import (
	"net/http"
)

type RequestLengthCheckPoint struct {
	CheckPoint
}

func (this *RequestLengthCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.ContentLength
	return
}
