package teaconfigs

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestServerConfig_NextBackend(t *testing.T) {
	a := assert.NewAssertion(t)

	s := NewServerConfig()
	s.Scheduling = &SchedulingConfig{
		Code: "random",
	}

	a.IsNil(s.NextBackend(nil))

	s.AddBackend(&BackendConfig{
		On:       true,
		IsBackup: false,
		IsDown:   false,
	})
	s.Validate()
	a.IsNotNil(s.NextBackend(nil))

	// backup
	s.Backends = []*BackendConfig{}
	s.AddBackend(&BackendConfig{
		Address:  ":80",
		On:       true,
		IsBackup: true,
		IsDown:   false,
		Weight:   10,
	})
	s.AddBackend(&BackendConfig{
		Address:  ":81",
		On:       true,
		IsBackup: false,
		IsDown:   false,
		Weight:   10,
	})
	s.Scheduling = &SchedulingConfig{
		Code: "roundRobin",
	}
	s.Validate()
	t.Log(s.schedulingObject.Summary())
	t.Log(s.NextBackend(nil))
	t.Log(s.NextBackend(nil))
	t.Log(s.NextBackend(nil))
	t.Log(s.NextBackend(nil))
}

func TestNewServerConfigFromId(t *testing.T) {
	t.Log(NewServerConfigFromId("123456"))
	t.Log(NewServerConfigFromId("defaultproxy"))
	t.Log(NewServerConfigFromId("XAPyHD1En4Z88yMc"))
}

func TestServerConfig_Encode(t *testing.T) {
	s := NewServerConfig()
	s.IgnoreHeaders = []string{"Server", "Content-Type"}
	s.AddBackend(NewBackendConfig())
	data, err := yaml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}
