package widgets

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/utils/string"
)

// 线图
type LineChart struct {
	Params []string `yaml:"params" json:"params"` // deprecated: v0.1.8 使用Lines代替
	Lines  []*Line  `yaml:"lines" json:"lines"`
}

// 添加线
func (this *LineChart) AddLine(line *Line) {
	this.Lines = append(this.Lines, line)
}

// 所有参数名
func (this *LineChart) AllParamNames() []string {
	result := []string{}
	for _, line := range this.Lines {
		teautils.ParseVariables(line.Param, func(varName string) (value string) {
			if !lists.ContainsString(result, varName) {
				result = append(result, varName)
			}
			return ""
		})
	}
	return result
}

// 转换为Javascript
func (this *LineChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	if len(this.Lines) == 0 {
		this.Lines = []*Line{}
	}

	// 兼容老的版本
	if len(this.Params) > 0 {
		for _, param := range this.Params {
			line := NewLine()
			line.Param = param
			this.AddLine(line)
		}
	}

	options["lines"] = this.Lines
	options["allParams"] = this.AllParamNames()
	return `
var chart = new charts.LineChart();
chart.options = ` + stringutil.JSONEncode(options) + `;

var query = NewQuery();
var ones = query.past(60, time.MINUTE).avg.apply(query, chart.options.allParams);

var lines = [];

chart.options.lines.$each(function (k, v) {
	var line = new charts.Line();
	if (v.color == null || v.color.length == 0) {
		line.color = (k < colors.ARRAY.length) ? colors.ARRAY[k] : null;
	} else {
		line.color = colors[v.color]
	}
	line.isFilled = v.isFilled;
	line.values = [];
	lines.push(line);
});

ones.$each(function (k, v) {
	chart.options.lines.$each(function (k, lineOption) {
		var value = values.valueOf(v.value, lineOption.param);
		lines[k].values.push(value);

		if (k == 0) {
			chart.addLabel(v.label);
		}
	});
});

chart.addLines(lines);
chart.render();
`, nil
}
