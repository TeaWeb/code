package checkpoints

import "net/http"

// Check Point
type CheckPointInterface interface {
	IsRequest() bool
	RequestValue(req *http.Request, param string) (value interface{}, err error)
	ResponseValue(req *http.Request, resp *http.Response, param string) (value interface{}, err error)
}
