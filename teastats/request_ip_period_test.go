package teastats

import (
	"fmt"
	"github.com/TeaWeb/code/tealogs"
	"testing"
	"time"
)

func TestRequestIPPeriodFilter_Start(t *testing.T) {
	queue := NewQueue()
	queue.Start("www")

	filter := new(RequestIPPeriodFilter)
	filter.Start(queue, "request.ip.day")
	t.Log(filter.rank.buffer)
	t.Log(filter.rank.min, filter.rank.minKey)

	for i := 0; i < 30; i ++ {
		accessLog := tealogs.NewAccessLog()
		accessLog.RemoteAddr = "192.168.1." + fmt.Sprintf("%d", i)
		filter.Filter(accessLog)
	}

	filter.Stop()

	time.Sleep(1 * time.Second)
	queue.Stop()
}
