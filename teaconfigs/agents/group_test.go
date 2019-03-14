package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSharedGroupConfig(t *testing.T) {
	config := SharedGroupConfig()
	logs.PrintAsJSON(config)
}

func TestSharedGroupConfig_Save(t *testing.T) {
	config := SharedGroupConfig()
	config.AddGroup(NewGroup("GROUP001"))
	err := config.Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSharedGroupConfig_Remove(t *testing.T) {
	config := SharedGroupConfig()
	config.RemoveGroup("0iqcKNj6zYCqfMoH")
	err := config.Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSharedGroupConfig2(t *testing.T) {
	config := SharedGroupConfig()
	//logs.PrintAsJSON(config, t)
	logs.PrintAsJSON(config.FindGroup(""), t)
}
