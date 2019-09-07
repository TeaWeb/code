package db

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
)

// MySQL配置
type MySQLConfig struct {
	DSN string `yaml:"dsn" json:"dsn"`
}

// 加载MySQL配置
func LoadMySQLConfig() (*MySQLConfig, error) {
	data, err := ioutil.ReadFile(Tea.ConfigFile("mysql.conf"))
	if err != nil {
		return nil, err
	}
	config := &MySQLConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
