package teastats

import (
	"testing"
	"time"
)

func TestKVStorage_Set(t *testing.T) {
	err := sharedKV.Set("hello", "value", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	value, err := sharedKV.Get("hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
	sharedKV.Close()
}

func TestKVStorage_Has(t *testing.T) {
	t.Log(sharedKV.Has("hello"))
}
