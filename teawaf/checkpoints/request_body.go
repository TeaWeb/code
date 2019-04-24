package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

// ${requestBody}
type RequestBodyCheckpoint struct {
	Checkpoint
}

func (this *RequestBodyCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	// TODO improve performance: ReadBody should be called once for one single request

	data, err := req.ReadBody(int64(32 * 1024 * 1024)) // read 32m bytes
	if err != nil {
		return "", err, nil
	}

	req.RestoreBody(data)

	return string(data), nil, nil
}

func (this *RequestBodyCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
