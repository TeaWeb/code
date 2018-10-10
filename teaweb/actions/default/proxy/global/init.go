package global

var proxyChanged = false

func NotifyChange() {
	proxyChanged = true
}

func FinishChange() {
	proxyChanged = false
}

func ProxyIsChanged() bool {
	return proxyChanged
}
