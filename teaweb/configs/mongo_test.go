package configs

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadMongoConfig(t *testing.T) {
	config, err := LoadMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(config, t)
}

func TestMongoConfig_Save(t *testing.T) {
	config, err := LoadMongoConfig()
	if err != nil {
		t.Fatal(err)
	}

	config.AccessLog = &MongoAccessLogConfig{
		CleanHour: 2,
		KeepDays:  3,
	}
	err = config.Save()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
