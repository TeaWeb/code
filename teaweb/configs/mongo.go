package configs

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
)

// MongoDB配置
type MongoConfig struct {
	URI       string                `yaml:"uri" json:"uri"`
	AccessLog *MongoAccessLogConfig `yaml:"accessLog" json:"accessLog"`
}

type MongoAccessLogConfig struct {
	CleanHour int `yaml:"cleanHour" json:"cleanHour"` // 清理时间，0-23
	KeepDays  int `yaml:"keepDays" json:"keepDays"`   // 保留挺熟
}

// 获取新对象
func NewMongoConfig() *MongoConfig {
	return &MongoConfig{}
}

// 加载MongoDB配置
func LoadMongoConfig() (*MongoConfig, error) {
	data, err := ioutil.ReadFile(Tea.ConfigFile("mongo.conf"))
	if err != nil {
		return nil, err
	}
	config := &MongoConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 保存
func (this *MongoConfig) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Tea.ConfigFile("mongo.conf"), data, 0666)
}
