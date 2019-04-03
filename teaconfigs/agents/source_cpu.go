package agents

import (
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"sync"
	"time"
)

// CPU Locker
var (
	cpuSourceLocker sync.Mutex
	cpuSourceTime   time.Time
	cpuSourceValue  interface{} = nil
)

// CPU使用量
type CPUSource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewCPUSource() *CPUSource {
	return &CPUSource{}
}

// 名称
func (this *CPUSource) Name() string {
	return "CPU"
}

// 代号
func (this *CPUSource) Code() string {
	return "cpu"
}

// 描述
func (this *CPUSource) Description() string {
	return "CPU使用量等信息"
}

// 执行
func (this *CPUSource) Execute(params map[string]string) (value interface{}, err error) {
	cpuSourceLocker.Lock()
	defer cpuSourceLocker.Unlock()

	if time.Since(cpuSourceTime).Seconds() < 1 && cpuSourceValue != nil {
		value = cpuSourceValue
		return
	}

	percents, err := cpu.Percent(0, true)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if len(percents) == 0 {
		value = map[string]interface{}{
			"avg": 0,
			"all": []float64{},
		}
		return
	}
	sum := float64(0)
	for _, percent := range percents {
		// 修复Windows上可能遇到的100%的Bug
		if runtime.GOOS == "windows" && percent > 99.9 {
			percent = 0
		}

		sum += percent
	}
	avg := sum / float64(len(percents))

	value = map[string]interface{}{
		"usage": map[string]interface{}{
			"avg": avg,
			"all": percents,
		},
	}

	cpuSourceTime = time.Now()
	cpuSourceValue = value

	return
}

// 表单信息
func (this *CPUSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *CPUSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "usage.avg",
			Description: "所有CPU平均使用量",
		},
		{
			Code:        "usage.all",
			Description: "每个CPU使用的量",
		},
		{
			Code:        "usage.all.$",
			Description: "单个CPU使用量，$表示0到N",
		},
	}
}

// 阈值
func (this *CPUSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${usage.avg}"
		t.Operator = ThresholdOperatorGte
		t.Value = "80"
		t.NoticeLevel = notices.NoticeLevelWarning
		result = append(result, t)
	}

	return result
}

// 图表
func (this *CPUSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "CPU使用量（%）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.LineChart();
chart.max = 100;

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].values.push(v.value.usage.avg);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.render();
`,
		}
		charts = append(charts, chart)
	}

	return charts
}
