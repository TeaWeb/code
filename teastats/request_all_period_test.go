package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"testing"
	"time"
)

func TestRequestAllPeriodFilter_Filter(t *testing.T) {
	queue := new(Queue)
	queue.Start("123456")

	filter := new(RequestAllPeriodFilter)
	filter.Start(queue, "request.all.minute")

	before := time.Now()
	for i := 0; i < 50000; i ++ {
		accessLog := &tealogs.AccessLog{}
		accessLog.Timestamp = time.Now().Unix()
		filter.Filter(accessLog)
		//time.Sleep(300 * time.Millisecond)
		accessLog.Timestamp = time.Now().Unix()
	}
	filter.Stop()
	t.Log(time.Since(before).Seconds(), "seconds")

	time.Sleep(1 * time.Second)
}
