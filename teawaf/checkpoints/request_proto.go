package checkpoints

import (
	"net/http"
)

type RequestProtoCheckPoint struct {
	CheckPoint
}

func (this *RequestProtoCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.Proto
	return
}
