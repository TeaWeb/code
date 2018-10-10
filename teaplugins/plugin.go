package teaplugins

type Plugin struct {
	Name      string
	Version   string
	Developer string
	Widgets   []*Widget
}

func NewPlugin() *Plugin {
	return &Plugin{
		Widgets: []*Widget{},
	}
}

func (this *Plugin) AddWidget(widget *Widget) {
	this.Widgets = append(this.Widgets, widget)
}
