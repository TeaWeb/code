package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestNetworkSource_Execute(t *testing.T) {
	source := NewNetworkSource()

	for i := 0; i < 2; i ++ {
		before := time.Now()
		value, err := source.Execute(nil)
		t.Log(time.Since(before).Seconds(), "s")
		if err != nil {
			t.Fatal(err)
		}
		logs.PrintAsJSON(value, t)
		time.Sleep(1 * time.Second)
	}
}
