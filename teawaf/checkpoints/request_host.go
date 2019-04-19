package checkpoints

import (
	"net/http"
)

type RequestHostCheckPoint struct {
	CheckPoint
}

func (this *RequestHostCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Host
	return
}
