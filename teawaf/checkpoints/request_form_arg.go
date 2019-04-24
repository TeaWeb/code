package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
	"net/url"
)

// ${requestForm.arg}
type RequestFormArgCheckpoint struct {
	Checkpoint
}

func (this *RequestFormArgCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	// TODO improve performance: ReadBody() should be called once for one single request
	data, err := req.ReadBody(32 * 1024 * 1024) // read 32m bytes
	if err != nil {
		return "", err, nil
	}

	values, _ := url.ParseQuery(string(data))

	req.RestoreBody(data)
	return values.Get(param), nil, nil
}

func (this *RequestFormArgCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
