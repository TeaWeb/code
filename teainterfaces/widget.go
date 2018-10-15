package teainterfaces

type WidgetGroup = uint8

const (
	WidgetGroupSystem   = WidgetGroup(1) // 系统信息
	WidgetGroupService  = WidgetGroup(2) // 服务
	WidgetGroupRealTime = WidgetGroup(3) // 即时
)

type WidgetInterface interface {
	Name() string
	Icon() []byte
	Title() string
	URL() string
	MoreURL() string
	TopBar() bool

	MenuBar() bool
	HelperBar() bool
	Dashboard() bool

	Group() WidgetGroup
	Charts() []interface{}

	OnReload() error
}
