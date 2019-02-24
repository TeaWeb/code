package teastats

import (
	"github.com/TeaWeb/code/tealogs"
)

// 服务队列配置
type ServerQueue struct {
	Queue   *Queue
	Filters map[string]FilterInterface // code => instance
}

func (this *ServerQueue) Stop() {
	for _, f := range this.Filters {
		f.Stop()
	}

	this.Queue.Stop()
	this.Queue = nil
	this.Filters = nil
}

func (this *ServerQueue) StartFilter(code string) {
	_, found := this.Filters[code]
	if found {
		return
	}

	instance := FindFilter(code)
	if instance == nil {
		return
	}

	this.Filters[code] = instance
	instance.Start(this.Queue, code)
}

func (this *ServerQueue) Filter(accessLog *tealogs.AccessLog) {
	for _, f := range this.Filters {
		f.Filter(accessLog)
	}
}
