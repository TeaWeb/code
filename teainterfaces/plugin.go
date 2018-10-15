package teainterfaces

type PluginInterface interface {
	Name() string
	Site() string // 网站
	Code() string
	Version() string
	Date() string
	Developer() string
	Description() string

	Widgets() []interface{}
	OnLoad() error
	OnReload() error
	OnStart() error
	OnStop() error
	OnUnload() error
}
