package agents

// 环境变量
type EnvVariable struct {
	Name  string `yaml:"name" json:"name"`   // 变量名
	Value string `yaml:"value" json:"value"` // 变量值
}
