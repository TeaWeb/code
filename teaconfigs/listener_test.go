package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestParseConfigs(t *testing.T) {
	configs, err := ParseConfigs()
	if err != nil {
		t.Error(err)
		return
	}
	for _, config := range configs {
		t.Log(config.Key, config.Address, config.Servers)
	}
}

func TestListenerConfig_FindNamedServer(t *testing.T) {
	a := assert.NewAssertion(t)

	listener := &ListenerConfig{}

	{
		server := NewServerConfig()
		server.AddName("h.com", "test.hello.com")
		listener.AddServer(server)
	}
	{
		server := NewServerConfig()
		server.AddName("hello.com")
		listener.AddServer(server)
	}

	{
		server := NewServerConfig()
		server.AddName("*.hello.com")
		listener.AddServer(server)
	}

	{
		result, _ := listener.FindNamedServer("hello.com")
		if result != nil {
			a.Log(result.Name)
		} else {
			a.Log("not found")
		}
	}

	{
		result, _ := listener.FindNamedServer("a.hello.com")
		if result != nil {
			a.Log(result.Name)
		} else {
			a.Log("not found")
		}
	}
}
