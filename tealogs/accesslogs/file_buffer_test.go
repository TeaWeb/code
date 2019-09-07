package accesslogs

import (
	"fmt"
	"github.com/TeaWeb/code/teautils"
	"github.com/go-acme/lego/log"
	"github.com/iwind/TeaGo/logs"
	"github.com/pquerna/ffjson/ffjson"
	"os"
	"testing"
	"time"
)

func TestFileBuffer_Write(t *testing.T) {
	buf := teautils.NewFileBuffer("hello.log")
	err := buf.Open()
	if err != nil {
		t.Fatal(err)
	}
	writeAccessLogToBuffer(buf, "hello")
	writeAccessLogToBuffer(buf, "world")

	max := 10000

	go func() {
		i := 0
		for {
			i++
			if i > max {
				break
			}
			if i%1000 == 0 {
				time.Sleep(1 * time.Second)
			}
			writeAccessLogToBuffer(buf, "Fine "+fmt.Sprintf("%d", i))
		}
	}()

	j := 0
	for {
		data, err := buf.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			time.Sleep(1 * time.Second)
			j++
			if j < max {
				continue
			} else {
				break
			}
		}
		log.Println("line:", len(data), string(data))
	}

	_ = os.Remove("hello.log")
}

func writeAccessLogToBuffer(buf *teautils.FileBuffer, path string) {
	accessLog := &AccessLog{
		RequestPath: path,
	}
	data, err := ffjson.Marshal(accessLog)
	if err != nil {
		logs.Error(err)
		return
	}
	err = buf.Write(data)
	if err != nil {
		logs.Error(err)
	}
}
