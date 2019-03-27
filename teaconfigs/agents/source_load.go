package agents

import (
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
)

// 负载
type LoadSource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewLoadSource() *LoadSource {
	return &LoadSource{}
}

// 名称
func (this *LoadSource) Name() string {
	return "负载"
}

// 代号
func (this *LoadSource) Code() string {
	return "load"
}

// 描述
func (this *LoadSource) Description() string {
	return "系统负载信息"
}

// 表单信息
func (this *LoadSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *LoadSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "load1",
			Description: "1分钟平均负载",
		},
		{
			Code:        "load5",
			Description: "5分钟平均负载",
		},
		{
			Code:        "load15",
			Description: "15分钟平均负载",
		},
	}
}

// 阈值
func (this *LoadSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${load5}"
		t.Value = "10"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGte
		result = append(result, t)
	}

	return result
}

// 图表
func (this *LoadSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "负载（Load）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.name = "1分钟";
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "5分钟";
	line.color = colors.BROWN;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "15分钟";
	line.color = colors.RED;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
}

var maxValue = 1;

ones.$each(function (k, v) {
	lines[0].values.push(v.value.load1);
	lines[1].values.push(v.value.load5);
	lines[2].values.push(v.value.load15);

	if (v.value.load1 > maxValue) {
		maxValue = Math.ceil(v.value.load1 / 2) * 2;
	}
	if (v.value.load5 > maxValue) {
		maxValue = Math.ceil(v.value.load5 / 2) * 2;
	}
	if (v.value.load15 > maxValue) {
		maxValue = Math.ceil(v.value.load15 / 2) * 2;
	}
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.max = maxValue;
chart.render();
`,
		}

		charts = append(charts, chart)
	}

	return charts
}
