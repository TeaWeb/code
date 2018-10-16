package teainterfaces

type PluginInterface interface {
	Name() string // 插件名
	Code() string
	Site() string // 网站
	Version() string
	Date() string // 发布日期
	Developer() string
	Description() string

	OnLoad() error
	OnReload() error
	OnStart() error
	OnStop() error
	OnUnload() error
}
