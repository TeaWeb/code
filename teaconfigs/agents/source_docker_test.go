package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestDockerSource_Execute(t *testing.T) {
	source := NewDockerSource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
