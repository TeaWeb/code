package checkpoints

import (
	"encoding/json"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
	"strings"
)

// ${requestJSON.arg}
type RequestJSONArgCheckpoint struct {
	Checkpoint
}

func (this *RequestJSONArgCheckpoint) RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	// TODO improve performance: ReadBody should be called once for one single request

	data, err := req.ReadBody(int64(32 * 1024 * 1024)) // read 32m bytes
	if err != nil {
		return "", err, nil
	}
	defer req.RestoreBody(data)

	var m interface{} = nil
	err = json.Unmarshal(data, &m)
	if err != nil || m == nil {
		return "", nil, err
	}

	value = teautils.Get(m, strings.Split(param, "."))
	if value != nil {
		return value, nil, err
	}
	return "", nil, nil
}

func (this *RequestJSONArgCheckpoint) ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param)
	}
	return
}
