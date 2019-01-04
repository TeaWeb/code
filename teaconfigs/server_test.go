package teaconfigs

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestServerConfig_NextBackend(t *testing.T) {
	a := assert.NewAssertion(t)

	s := NewServerConfig()
	s.Scheduling = &SchedulingConfig{
		Code: "random",
	}

	a.IsNil(s.NextBackend(maps.Map{}))

	s.AddBackend(&BackendConfig{
		On:       true,
		IsBackup: false,
		IsDown:   false,
	})
	s.Validate()
	a.IsNotNil(s.NextBackend(maps.Map{}))

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
	t.Log(s.NextBackend(maps.Map{}))
	t.Log(s.NextBackend(maps.Map{}))
	t.Log(s.NextBackend(maps.Map{}))
	t.Log(s.NextBackend(maps.Map{}))
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
