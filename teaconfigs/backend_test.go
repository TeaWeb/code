package teaconfigs

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/assert"
	"sync"
	"testing"
	"time"
)

func TestBackendConfig(t *testing.T) {
	yamlData, err := yaml.Marshal(new(BackendConfig))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))
}

func TestBackendConfig_IncreaseConn(t *testing.T) {
	backend := NewBackendConfig()
	count := 10000
	wg := sync.WaitGroup{}
	wg.Add(count)
	before := time.Now()
	for i := 0; i < count; i ++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j ++ {
				backend.IncreaseConn()
			}
		}()
	}
	wg.Wait()
	t.Log(float64(count)/time.Since(before).Seconds(), "qps")
	t.Log("result:", backend.CurrentConns)
}

func TestBackendConfig_DecreaseConn(t *testing.T) {
	backend := NewBackendConfig()
	backend.CurrentConns = 10000000

	count := 10000
	wg := sync.WaitGroup{}
	wg.Add(count)
	before := time.Now()
	for i := 0; i < count; i ++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j ++ {
				backend.DecreaseConn()
			}
		}()
	}
	wg.Wait()
	t.Log(float64(count)/time.Since(before).Seconds(), "qps")
	t.Log("result:", backend.CurrentConns)
}

func TestBackendConfig_RequestPath(t *testing.T) {
	a := assert.NewAssertion(t)
	{
		backend := NewBackendConfig()
		backend.Validate()
		a.IsFalse(backend.HasRequestURI())
	}

	{
		backend := NewBackendConfig()
		backend.RequestURI = "${requestURI}"
		backend.Validate()
		a.IsFalse(backend.HasRequestURI())
	}

	{
		backend := NewBackendConfig()
		backend.RequestURI = "/hello${requestURI}"
		backend.Validate()
		a.IsTrue(backend.HasRequestURI())
		a.IsTrue(backend.RequestPath() == "/hello${requestURI}")
		a.IsTrue(backend.RequestArgs() == "")
	}

	{
		backend := NewBackendConfig()
		backend.RequestURI = "/hello${requestURI}?name=value"
		backend.Validate()
		a.IsTrue(backend.HasRequestURI())
		a.IsTrue(backend.RequestPath() == "/hello${requestURI}")
		a.IsTrue(backend.RequestArgs() == "name=value")
	}
}
