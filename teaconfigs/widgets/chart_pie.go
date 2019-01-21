package widgets

import (
	"github.com/iwind/TeaGo/utils/string"
)

// Pie
type PieChart struct {
	Param string
	Limit int
}

func (this *PieChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	options["limit"] = this.Limit
	options["param"] = this.Param
	return `
var chart = new charts.PieChart();
chart.options = ` + stringutil.JSONEncode(options) + `;

var query = new values.Query();
if (chart.options.limit <= 0) {
	query.limit(100);
} else {
	query.limit(chart.options.limit);
}
var ones = query.desc().cache(60).findAll();
chart.values = [];
chart.labels = [];
ones.$each(function (k, v) {
	var value = values.valueOf(v.value, chart.options.param);
	var index = chart.labels.$indexOf(value);
	if (index == -1) {
		chart.values.push(1);
		chart.labels.push(value);
	} else {
		chart.values[index] ++;
	}
});

chart.render();
`, nil
}
