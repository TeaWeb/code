package teaproxy

import (
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

// 调用代理
func (this *Request) callProxy(writer *ResponseWriter) error {
	options := maps.Map{
		"request":   this.raw,
		"formatter": this.Format,
	}
	backend := this.proxy.NextBackend(options)

	responseCallback := options.Get("responseCallback")
	if responseCallback != nil {
		f, ok := responseCallback.(func(http.ResponseWriter))
		if ok {
			this.responseCallback = f
		}
	}

	this.backend = backend
	return this.callBackend(writer)
}
