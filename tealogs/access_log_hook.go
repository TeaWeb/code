package tealogs

// 外部Hook
var accessLogHooks = []*AccessLogHook{}

// 请求Hook定义
type AccessLogHook struct {
	Process func(accessLog *AccessLog) (goNext bool)
}

// 添加Hook
func AddAccessLogHook(hook *AccessLogHook) {
	accessLogHooks = append(accessLogHooks, hook)
}

// 执行Filter
func CallAccessLogHooks(accessLog *AccessLog) {
	if len(accessLogHooks) == 0 {
		return
	}
	for _, hook := range accessLogHooks {
		goNext := hook.Process(accessLog)
		if !goNext {
			break
		}
	}
}
