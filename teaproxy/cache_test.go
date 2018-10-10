package teaproxy

import "testing"

func TestNewFixedCache(t *testing.T) {
	cache := NewFixedCache()
	t.Log(cache.maxMemory / 1024 / 1024 / 1024)
}

func TestFixedCache_Add(t *testing.T) {
	cache := NewFixedCache()
	cache.Add("hello", "world")
	cache.Add("name", "liu")
	cache.Add("age", "20")
	cache.Add("var1", "20")
	cache.Add("var2", "20")
	cache.Add("var3", "20")

	t.Log(cache.Get("hello"))
	t.Log(cache.Get("name"))
	t.Log(cache.Get("age"))

	t.Log("=====")
	cache.Trim()

	t.Log(cache.Get("hello"))
	t.Log(cache.Get("name"))
	t.Log(cache.Get("age"))
	t.Log(cache.Get("var1"))
	t.Log(cache.Get("var2"))
	t.Log(cache.Get("var3"))
}