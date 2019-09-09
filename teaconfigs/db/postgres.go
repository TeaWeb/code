package db

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
)

// PostgreSQL配置
type PostgresConfig struct {
	DSN string `yaml:"dsn" json:"dsn"`
}

// 加载PostgreSQL配置
func LoadPostgresConfig() (*PostgresConfig, error) {
	data, err := ioutil.ReadFile(Tea.ConfigFile("postgres.conf"))
	if err != nil {
		return nil, err
	}
	config := &PostgresConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
