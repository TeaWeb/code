package teautils

import (
	"fmt"
	"github.com/TeaWeb/code/teatesting"
	"os"
	"testing"
	"time"
)

func TestStatFile(t *testing.T) {
	t.Log(StatFile("os.go"))
	t.Log(StatFile("os.go"))

	if !teatesting.IsGlobal() {
		time.Sleep(10 * time.Second)
	}

	t.Log(StatFile("os.go"))
}

func TestStatFile2(t *testing.T) {
	_, _ = StatFile("os.go")
	_, _ = StatFile("os_test.go")
	_, _ = StatFile("service.go")
	_, _ = StatFile("string.go")
	_, _ = StatFile("string_test.go")

	time.Sleep(1 * time.Second)

	_, _ = StatFile("string_test.go")
}

func TestStatManyFiles(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	_, err := os.Stat("/tmp")
	if err != nil {
		return
	}
	for i := 0; i < 10000; i++ {
		file := "/tmp/stat." + fmt.Sprintf("%d", i) + ".log"
		fp, err := os.Create(file)
		if err != nil {
			t.Fatal(err)
		} else {
			_ = fp.Close()
		}

		_, _ = StatFile(file)
	}

	time.Sleep(30 * time.Second)

	for i := 0; i < 10000; i++ {
		file := "/tmp/stat." + fmt.Sprintf("%d", i) + ".log"
		_ = os.Remove(file)
	}
}

func BenchmarkStatFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StatFile("os.go")
	}
}

func BenchmarkStatFile_Raw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = os.Stat("os.go")
	}
}
