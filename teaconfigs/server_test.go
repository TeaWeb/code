package teaconfigs

import (
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

	s.AddBackend(&ServerBackendConfig{
		On:       true,
		IsBackup: false,
		IsDown:   false,
	})
	s.Validate()
	a.IsNotNil(s.NextBackend(maps.Map{}))

	// backup
	s.Backends = []*ServerBackendConfig{}
	s.AddBackend(&ServerBackendConfig{
		Address:  ":80",
		On:       true,
		IsBackup: true,
		IsDown:   false,
		Weight:   10,
	})
	s.AddBackend(&ServerBackendConfig{
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
