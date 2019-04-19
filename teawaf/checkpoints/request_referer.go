package checkpoints

import (
	"net/http"
)

type RequestRefererCheckPoint struct {
	CheckPoint
}

func (this *RequestRefererCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Referer()
	return
}
