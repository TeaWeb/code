package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"testing"
	"time"
)

func TestIPAllPeriodFilter_Start(t *testing.T) {
	queue := new(Queue)
	queue.ServerId = "123456"

	filter := new(IPAllPeriodFilter)
	filter.Start(queue, "ip.all.minute")

	accessLog := &tealogs.AccessLog{}
	accessLog.Timestamp = time.Now().Unix()
	accessLog.RemoteAddr = "127.0.0.1:1234"
	filter.Filter(accessLog)
	t.Log(filter.values)
}
