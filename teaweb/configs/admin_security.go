package configs

// 安全设置定义
type AdminSecurity struct {
	Allow      []string `yaml:"allow" json:"allow"`
	Deny       []string `yaml:"deny" json:"deny"`
	Secret     string   `yaml:"secret" json:"secret"`
	IsDisabled bool     `yaml:"isDisabled" json:"isDisabled"` // 是否禁用
}
