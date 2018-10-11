package configs

type AdminUser struct {
	Username string   `yaml:"username" json:"username"`
	Password string   `yaml:"password" json:"password"`
	Role     []string `yaml:"role" json:"role"`
}
