package teaplugins

const (
	WidgetGroupSystem   = WidgetGroup(1) // 系统信息
	WidgetGroupService  = WidgetGroup(2) // 服务
	WidgetGroupRealTime = WidgetGroup(3) // 即时
)

type WidgetGroup uint8

type Group struct {
	Id      WidgetGroup `json:"id"`
	Name    string      `json:"name"`
	Widgets []*Widget   `json:"widgets"`
}

func DashboardGroups() []*Group {
	return []*Group{
		{
			Id:      WidgetGroupRealTime,
			Name:    "实时信息",
			Widgets: DashboardWidgets(WidgetGroupRealTime),
		},
		{
			Id:      WidgetGroupSystem,
			Name:    "系统信息",
			Widgets: DashboardWidgets(WidgetGroupSystem),
		},

		{
			Id:      WidgetGroupService,
			Name:    "服务",
			Widgets: DashboardWidgets(WidgetGroupService),
		},
	}
}

func (this *Group) Reload() {
	for _, widget := range this.Widgets {
		widget.Reload()
	}
}

func (this *Group) ForceReload() {
	for _, widget := range this.Widgets {
		widget.ForceReload()
	}
}
