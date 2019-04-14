package teacache

import (
	"testing"
	"time"
)

func TestLevelDBManager_Write(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": "cache",
	})
	m.Life = 30 * time.Second
	t.Log(m.Write("hello", []byte("world123")))
}

func TestLevelDBManager_Read(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": "cache",
	})
	data, err := m.Read("hello")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestLevelDBManager_Stat(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": "cache",
	})
	t.Log(m.Stat())
}

func TestLevelDBManager_CleanExpired(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": "cache",
	})
	m.Life = 30 * time.Second
	m.CleanExpired()
}

func TestLevelDBManager_Clean(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": "cache",
	})
	m.Life = 30 * time.Second
	m.Clean()
}
