package checkpoints

import (
	"net/http"
)

type RequestCookieCheckpoint struct {
	Checkpoint
}

func (this *RequestCookieCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	cookie, err := req.Cookie(param)
	if err != nil {
		value = ""
		return
	}

	value = cookie.Value
	return
}
