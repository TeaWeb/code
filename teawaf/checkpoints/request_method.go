package checkpoints

import (
	"net/http"
)

type RequestMethodCheckPoint struct {
	CheckPoint
}

func (this *RequestMethodCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Method
	return
}
