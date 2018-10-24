package teaplugins

import (
	"github.com/TeaWeb/code/teaapps"
	"github.com/iwind/TeaGo/utils/string"
)

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
	Apps        []*teaapps.App
}

func NewPlugin() *Plugin {
	return &Plugin{
		Widgets: []*Widget{},
	}
}

func (this *Plugin) AddWidget(widget *Widget) {
	if len(widget.Id) == 0 {
		widget.Id = stringutil.Rand(16)
	}

	this.Widgets = append(this.Widgets, widget)
}

func (this *Plugin) ResetApps() {
	this.Apps = []*teaapps.App{}
}

func (this *Plugin) AddApp(app *teaapps.App) {
	if len(app.Id) == 0 {
		app.Id = stringutil.Rand(16)
	}
	this.Apps = append(this.Apps, app)
}

func (this *Plugin) AppWithId(appId string) *teaapps.App {
	for _, p := range this.Apps {
		if p.Id == appId {
			return p
		}
	}
	return nil
}

func (this *Plugin) InterfaceNames() []string {
	names := []string{}
	if len(this.Widgets) > 0 {
		names = append(names, "widget")
	}
	if len(this.Apps) > 0 {
		names = append(names, "app")
	}
	return names
}
