package configs

import (
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v2"
)

type AdminConfig struct {
	Security struct {
		Allow  []string `yaml:"allow"`
		Deny   []string `yaml:"deny"`
		Secret string   `yaml:"secret"`
	} `yaml:"security"`

	Roles []struct {
		Name  string   `yaml:"name"`
		Grant []string `yaml:"grant"`
	} `yaml:"roles"`

	Users []struct {
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
		Role     []string `yaml:"role"`
	} `yaml:"users"`
}

var adminConfig *AdminConfig

func SharedAdminConfig() *AdminConfig {
	if adminConfig != nil {
		return adminConfig
	}

	adminConfig = &AdminConfig{}

	configFile := Tea.ConfigFile("admin.conf")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		logs.Error(err)
		return adminConfig
	}

	err = yaml.Unmarshal(data, adminConfig)
	if err != nil {
		logs.Error(err)
		return adminConfig
	}

	return adminConfig
}
