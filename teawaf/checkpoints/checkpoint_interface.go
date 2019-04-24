package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

// Check Point
type CheckpointInterface interface {
	IsRequest() bool
	RequestValue(req *requests.Request, param string) (value interface{}, sysErr error, userErr error)
	ResponseValue(req *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error)
	ParamOptions() *ParamOptions
}
