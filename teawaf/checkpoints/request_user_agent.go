package checkpoints

import (
	"net/http"
)

type RequestUserAgentCheckpoint struct {
	Checkpoint
}

func (this *RequestUserAgentCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = req.UserAgent()
	return
}
