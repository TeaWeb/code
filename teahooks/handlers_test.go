package teahooks

import (
	"testing"
)

func TestHandlers(t *testing.T) {
	On(EventReload, func() {
		t.Log("reload1")
	})
	On(EventReload, func() {
		t.Log("reload2")
	})
	Call(EventReload)
}
