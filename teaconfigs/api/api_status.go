package api

// API状态常量
const (
	APIStatusTypeNormal  = "normal"
	APIStatusTypeSuccess = "success"
	APIStatusTypeWarning = "warning"
	APIStatusTypeFailure = "failure"
	APIStatusTypeError   = "error"
)

// API状态定义
type APIStatus struct {
	Code        uint     `yaml:"code" json:"code"`               // 代码
	Description string   `yaml:"description" json:"description"` // 描述
	Groups      []string `yaml:"groups" json:"groups"`           // 分组
	Versions    []string `yaml:"versions" json:"versions"`       // 版本
	Type        string   `yaml:"type" json:"type"`               // 类型
}
