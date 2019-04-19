package checkpoints

import (
	"net/http"
)

type RequestUserAgentCheckPoint struct {
	CheckPoint
}

func (this *RequestUserAgentCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.UserAgent()
	return
}
