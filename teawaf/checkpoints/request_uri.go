package checkpoints

import (
	"net/http"
)

type RequestURICheckPoint struct {
	CheckPoint
}

func (this *RequestURICheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.RequestURI
	return
}
