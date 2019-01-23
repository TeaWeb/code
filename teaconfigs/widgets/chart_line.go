package widgets

import "github.com/iwind/TeaGo/utils/string"

// 线图
type LineChart struct {
	Params []string
	Limit  int
}

func (this *LineChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	// 防止出现null
	if len(this.Params) == 0 {
		this.Params = []string{}
	}

	options["limit"] = this.Limit
	options["params"] = this.Params
	return `
var chart = new charts.LineChart();
chart.options = ` + stringutil.JSONEncode(options) + `;

var query = new values.Query();
if (chart.options.limit <= 0) {
	query.limit(10);
} else {
	query.limit(chart.options.limit);
}
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];
chart.options.params.$each(function (k, v) {
	var line = new charts.Line();
	line.color = (k < colors.ARRAY.length) ? colors.ARRAY[k] : null;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
});
ones.$each(function (k, v) {
	chart.options.params.$each(function (k, param) {
		var value = values.valueOf(v.value, param);
		lines[k].values.push(value);
	});
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.render();
`, nil
}
