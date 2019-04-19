package checkpoints

import (
	"net/http"
)

type RequestSchemeCheckPoint struct {
	CheckPoint
}

func (this *RequestSchemeCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.URL.Scheme
	return
}
