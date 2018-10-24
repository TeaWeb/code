package teaapps

// 对App的操作定义
type Operation struct {
	Id            string // 唯一ID，通常系统会自动生成
	Code          string // 代号
	Name          string // 名称
	IsEnabled     bool   // 是否启用
	ShouldConfirm bool   // 操作是否确认
	ConfirmText   string // 操作确认文字

	onRunFunc func() error // 回调函数
}

// 设置运行时回调函数
func (this *Operation) OnRun(f func() error) {
	this.onRunFunc = f
}

// 运行此操作
func (this *Operation) Run() error {
	if this.onRunFunc != nil {
		return this.onRunFunc()
	}
	return nil
}
