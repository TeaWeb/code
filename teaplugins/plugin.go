package teaplugins

type Plugin struct {
	IsExternal bool // 是否第三方开发的

	Name        string    // 名称
	Code        string    // 代号
	Version     string    // 版本
	Date        string    // 发布日期
	Site        string    // 网站链接
	Developer   string    // 开发者
	Description string    // 插件简介
	Widgets     []*Widget // 小组件

	interfaceNames []string
}

func NewPlugin() *Plugin {
	return &Plugin{
		Widgets: []*Widget{},
	}
}

func (this *Plugin) AddWidget(widget *Widget) {
	this.Widgets = append(this.Widgets, widget)
}

func (this *Plugin) AddInterfaceName(interfaceName string) {
	this.interfaceNames = append(this.interfaceNames, interfaceName)
}

func (this *Plugin) InterfaceNames() []string {
	return this.interfaceNames
}
