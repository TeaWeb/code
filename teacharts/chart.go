package teacharts

import "sync"

type ChartInterface interface {
	UniqueId() string
	SetUniqueId(id string)
	Reload()
}

type Chart struct {
	Id     string `json:"id"` // 用来标记图标的唯一性，可以不填，系统会自动生成
	Type   string `json:"type"`
	Name   string `json:"name"`
	Detail string `json:"detail"`

	onReloadFuncs []func()
	locker        sync.Mutex
}

func (this *Chart) OnReload(f func()) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.onReloadFuncs = append(this.onReloadFuncs, f)
}

func (this *Chart) Reload() {
	if len(this.onReloadFuncs) == 0 {
		return
	}
	go func() {
		for _, f := range this.onReloadFuncs {
			f()
		}
	}()
}
