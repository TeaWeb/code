package checkpoints

import (
	"net/http"
	"strings"
)

type RequestHeaderCheckPoint struct {
	CheckPoint
}

func (this *RequestHeaderCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	v, found := req.Header[param]
	if !found {
		value = ""
		return
	}
	value = strings.Join(v, ";")
	return
}
