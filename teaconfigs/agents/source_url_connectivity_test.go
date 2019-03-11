package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestURLConnectivitySource_Execute(t *testing.T) {
	source := NewURLConnectivitySource()
	source.URL = "https://baidu.com/"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
