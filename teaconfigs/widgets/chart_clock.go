package widgets

import "github.com/iwind/TeaGo/utils/string"

// 时钟
type ClockChart struct {
}

func (this *ClockChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	return `
var chart = new charts.Clock();
chart.options = ` + stringutil.JSONEncode(options) + `;

var result = new values.Query().latest(1);
if (result.length > 0) {
	chart.timestamp = result[0].timestamp + context.item.interval;
}

chart.render();
`, nil
}
