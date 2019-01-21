package widgets

import (
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestJavascriptChart_AsJavascript(t *testing.T) {
	c := new(JavascriptChart)
	c.Code = `
var chart = new charts.HTMLChart();
chart.render();
`
	t.Log(c.AsJavascript(maps.Map{
		"name":    "Hello,World",
		"columns": 2,
	}))
}
