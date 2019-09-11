package teautils

import "testing"

func TestBytesToString(t *testing.T) {
	t.Log(BytesToString([]byte("Hello,World")))
}

func BenchmarkBytesToString(b *testing.B) {
	data := []byte("Hello,World")
	for i := 0; i < b.N; i++ {
		_ = BytesToString(data)
	}
}

func BenchmarkBytesToString2(b *testing.B) {
	data := []byte("Hello,World")
	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}

func TestFormatAddress(t *testing.T) {
	t.Log(FormatAddress("127.0.0.1:1234"))
	t.Log(FormatAddress("127.0.0.1 : 1234"))
	t.Log(FormatAddress("127.0.0.1：1234"))
}