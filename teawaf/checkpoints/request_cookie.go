package checkpoints

import (
	"net/http"
)

type RequestCookieCheckPoint struct {
	CheckPoint
}

func (this *RequestCookieCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	cookie, err := req.Cookie(param)
	if err != nil {
		value = ""
		return
	}

	value = cookie.Value
	return
}
