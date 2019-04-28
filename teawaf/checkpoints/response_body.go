package checkpoints

import (
	"bytes"
	"github.com/TeaWeb/code/teawaf/requests"
	"io/ioutil"
)

// ${responseBody}
type ResponseBodyCheckpoint struct {
	Checkpoint
}

func (this *ResponseBodyCheckpoint) IsRequest() bool {
	return false
}

func (this *ResponseBodyCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = ""
	return
}

func (this *ResponseBodyCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = ""
	if resp != nil && resp.Body != nil {
		if len(resp.BodyData) > 0 {
			value = string(resp.BodyData)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			sysErr = err
			return
		}
		resp.BodyData = body
		resp.Body.Close()
		value = string(body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	return
}
