package teautils

import (
	"fmt"
	"testing"
	"time"
)

func TestParseVariables(t *testing.T) {
	v := ParseVariables("hello, ${name}", func(s string) string {
		return "Lu"
	})
	t.Log(v)
}

func TestParseVariablesPerformance(t *testing.T) {
	count := 10000
	before := time.Now()
	for i := 0; i < count; i ++ {
		ParseVariables("hello, ${name} "+fmt.Sprintf("%d", i), func(s string) string {
			return "Lu"
		})
	}
	cost := time.Since(before).Seconds()
	t.Log(float64(count)/cost, "qps")
}
