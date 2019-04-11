package teacache

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestCacheMemoryConfig(t *testing.T) {
	m := NewMemoryManager()
	m.Capacity = 1024
	m.Life = 30 * time.Second
	err := m.Write("/hello", []byte("Hello, World"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("memory:", m.memory, "bytes")

	a := assert.NewAssertion(t).Quiet()

	_, err = m.Read("hello")
	a.IsNotNil(err)

	data, err := m.Read("/hello")
	if err != nil {
		t.Fatal(err)
	}
	a.Equals(string(data), "Hello, World")

	t.Log(string(data))
}
