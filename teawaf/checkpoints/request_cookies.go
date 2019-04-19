package checkpoints

import (
	"net/http"
)

type RequestCookiesCheckPoint struct {
	CheckPoint
}

func (this *RequestCookiesCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return
}
