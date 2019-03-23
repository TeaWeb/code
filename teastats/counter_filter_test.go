package teastats

import (
	"testing"
	"time"
)

func TestCounterFilter_EncodeParams(t *testing.T) {
	filter := &CounterFilter{}
	{
		b := filter.encodeParams(map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
		})
		t.Log(b.String())
	}

	{
		before := time.Now()
		count := 10000
		for i := 0; i < count; i ++ {
			filter.encodeParams(map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
			})
		}
		t.Log(float64(count)/time.Since(before).Seconds(), "qps")
	}
}
