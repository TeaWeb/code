package configs

type AdminRole struct {
	Name  string   `yaml:"name" json:"name"`
	Grant []string `yaml:"grant" json:"grant"`
}
