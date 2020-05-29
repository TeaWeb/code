package utils

import (
	"regexp"
	"strings"
	"testing"
)

func TestMatchStringCache(t *testing.T) {
	regex := regexp.MustCompile(`\d+`)
	t.Log(MatchStringCache(regex, "123"))
	t.Log(MatchStringCache(regex, "123"))
	t.Log(MatchStringCache(regex, "123"))
}

func TestMatchBytesCache(t *testing.T) {
	regex := regexp.MustCompile(`\d+`)
	t.Log(MatchBytesCache(regex, []byte("123")))
	t.Log(MatchBytesCache(regex, []byte("123")))
	t.Log(MatchBytesCache(regex, []byte("123")))
}

func BenchmarkMatchStringCache(b *testing.B) {
	data := strings.Repeat("HELLO", 512)
	regex := regexp.MustCompile(`(?iU)\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\b`)

	for i := 0; i < b.N; i++ {
		_ = MatchStringCache(regex, data)
	}
}

func BenchmarkMatchStringCache_WithoutCache(b *testing.B) {
	data := strings.Repeat("HELLO", 512)
	regex := regexp.MustCompile(`(?iU)\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\b`)

	for i := 0; i < b.N; i++ {
		_ = regex.MatchString(data)
	}
}
