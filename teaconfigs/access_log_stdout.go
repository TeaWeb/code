package teaconfigs

// 日志stdout配置
type AccessLogStdoutConfig struct {
	Format string `yaml:"format"`
	Buffer string `yaml:"buffer"` // @TODO
	Flush  string `yaml:"flush"`  // @TODO
}
