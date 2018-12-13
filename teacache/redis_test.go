package teacache

import (
	"testing"
	"time"
)

func TestRedisManager(t *testing.T) {
	manager := NewRedisManager()
	manager.Life = 30 * time.Second
	manager.SetOptions(map[string]interface{}{
		"host": "127.0.0.1",
	})

	t.Log(manager.Write("hello", []byte("world")))
	r, err := manager.Read("hello")
	if err != nil {
		t.Fatal("err:", err)
	} else {
		t.Log("read:", string(r))
	}
}
