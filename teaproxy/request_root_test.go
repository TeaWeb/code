package teaproxy

import (
	"fmt"
	"github.com/dchest/siphash"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"runtime"
	"strconv"
	"testing"
)

func TestFileEtag(t *testing.T) {
	etag := stringutil.Md5(fmt.Sprintf("%d,%d", 1563192836000, 1024))
	t.Log(etag)
}

func TestFileEtag_hash(t *testing.T) {
	etag := siphash.Hash(0, 0, []byte("123.txt"+strconv.FormatInt(1563192836000, 10)+strconv.FormatInt(1024, 10)))
	t.Log(fmt.Sprintf("%0x", etag))
}

func TestFileEtag_str(t *testing.T) {
	etag := siphash.Hash(0, 0, []byte("123.txt"+strconv.FormatInt(1563192836000, 10)+strconv.FormatInt(1024, 10)))
	t.Log(fmt.Sprintf("%0x", etag))
}

func BenchmarkFileEtag(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		_ = stringutil.Md5("123.txt" + fmt.Sprintf("%d,%d", 1563192836000, 1024))
	}
}

func BenchmarkFileEtag_Hash(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%0x", siphash.Hash(0, 0, []byte("123.txt"+strconv.FormatInt(1563192836000, 10)+strconv.FormatInt(1024, 10))))
	}
}
