package configs

type AdminSecurity struct {
	Allow  []string `yaml:"allow" json:"allow"`
	Deny   []string `yaml:"deny" json:"deny"`
	Secret string   `yaml:"secret" json:"secret"`
}
