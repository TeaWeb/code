package checkpoints

import (
	"net/http"
	"net/url"
	"strings"
)

type RequestCookiesCheckPoint struct {
	CheckPoint
}

func (this *RequestCookiesCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	var cookies = []string{}
	for _, cookie := range req.Cookies() {
		cookies = append(cookies, url.QueryEscape(cookie.Name)+"="+url.QueryEscape(cookie.Value))
	}
	value = strings.Join(cookies, "&")
	return
}
