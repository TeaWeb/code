package teacache

import (
	"github.com/TeaWeb/code/teaproxy"
)

func init() {
	hook := &teaproxy.RequestHook{
		BeforeRequest: ProcessBeforeRequest,
		AfterRequest:  ProcessAfterRequest,
	}
	teaproxy.AddRequestHook(hook)
}
