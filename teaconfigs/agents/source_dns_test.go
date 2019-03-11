package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestDNSSource_Execute(t *testing.T) {
	source := NewDNSSource()
	source.Domain = "teaos.cn"
	source.Type = "A"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}

func TestDNSSource_Execute_CHANGE(t *testing.T) {
	source := NewDNSSource()
	source.Domain = "teaos.cn"
	source.Type = "CHANGE"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}

func TestDNSSource_Execute_MX(t *testing.T) {
	source := NewDNSSource()
	source.Domain = "teaos.cn"
	source.Type = "MX"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}

func TestDNSSource_Execute_NS(t *testing.T) {
	source := NewDNSSource()
	source.Domain = "teaos.cn"
	source.Type = "NS"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}

func TestDNSSource_Execute_TXT(t *testing.T) {
	source := NewDNSSource()
	source.Domain = "teaos.cn"
	source.Type = "TXT"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}
