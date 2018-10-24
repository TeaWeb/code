package teaapps

type LogType = string
type LogFormat = string

// 日志类型
const (
	LogTypeNormal = "normal"
	LogTypeAccess = "access"
	LogTypeError  = "error"
)

// 日志文件格式
const (
	LogFormatPlain = "plain"
	LogFormatXML   = "xml"  // @TODO 暂不支持
	LogFormatGzip  = "gzip" // @TODO 暂不支持
	LogFormatJSON  = "json" // @TODO 暂不支持
	LogFormatSQL   = "sql"  // @TODO 暂不支持
)

// 日志定义
type Log struct {
	Id     string    // 唯一ID，通常系统会自动生成
	Name   string    // 日志名
	Path   string    // 文件路径
	Type   LogType   // 类型
	Format LogFormat // 文件格式
}
