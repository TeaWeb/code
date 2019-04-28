package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"strings"
)

type RequestHeadersCheckpoint struct {
	Checkpoint
}

func (this *RequestHeadersCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	var headers = []string{}
	for k, v := range req.Header {
		for _, subV := range v {
			headers = append(headers, k+": "+subV)
		}
	}
	value = strings.Join(headers, "\n")
	return
}

func (this *RequestHeadersCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
