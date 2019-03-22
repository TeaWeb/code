package shared

// 键值对
type Pair struct {
	Name  string `yaml:"name" json:"name"`   // 变量名
	Value string `yaml:"value" json:"value"` // 变量值
}
