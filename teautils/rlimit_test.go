package teautils

import (
	"os"
	"testing"
)

func TestRlimit(t *testing.T) {
	SetRLimit(20480)
	for i := 0; i < 10240; i++ {
		_, err := os.Open("ulimit_test.go")
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("OK")
}

func TestSetSuitableRLimit(t *testing.T) {
	SetSuitableRLimit()
}
