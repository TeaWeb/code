package teaconfigs

// 日志文件配置
type AccessLogFileConfig struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
	Buffer string `yaml:"buffer"` // @TODO
	Flush  string `yaml:"flush"`  // @TODO
}
