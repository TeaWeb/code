package agents

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// App进程
type AppProcessesSource struct {
	Source `yaml:",inline"`

	CmdlineKeyword string `yaml:"cmdlineKeyword" json:"cmdlineKeyword"` // 命令行匹配关键词
}

// 获取新对象
func NewAppProcessesSource() *AppProcessesSource {
	return &AppProcessesSource{}
}

// 名称
func (this *AppProcessesSource) Name() string {
	return "App进程数"
}

// 代号
func (this *AppProcessesSource) Code() string {
	return "app.processes"
}

// 描述
func (this *AppProcessesSource) Description() string {
	return "获取某个App的进程数，依赖系统安装ps、grep命令"
}

// 执行
func (this *AppProcessesSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.CmdlineKeyword) == 0 {
		value = map[string]interface{}{
			"processes": 0,
		}
		err = errors.New("缺少命令行匹配关键词")
		return
	}

	exec := teautils.NewCommandExecutor()
	exec.Add("ps", "ax", "-o", "pid,command")
	exec.Add("grep", this.CmdlineKeyword)
	exec.Add("grep", "-v", " grep ")
	exec.Add("wc", "-l")
	output, err := exec.Run()
	if err != nil {
		return maps.Map{
			"processes": 0,
		}, err
	}

	return maps.Map{
		"processes": types.Int(output),
	}, nil
}

// 表单信息
func (this *AppProcessesSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewTextField("命令行匹配关键词", "Cmdline")
			field.IsRequired = true
			field.Code = "cmdlineKeyword"
			field.Comment = "比如mysql、mongod之类的能匹配你要监控的进程命令行的关键词"
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入命令行匹配关键词");
}
`
			group.Add(field)
		}
	}
	return form
}

// 变量
func (this *AppProcessesSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "processes",
			Description: "进程数",
		},
	}
}

// 阈值
func (this *AppProcessesSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${processes}"
		t.Operator = ThresholdOperatorEq
		t.Value = "0"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "App未启动进程"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *AppProcessesSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "App进程数"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var line = new charts.Line();
line.color = colors.ARRAY[0];
line.isFilled = true;
line.values = [];

ones.$each(function (k, v) {
	line.values.push(v.value.processes);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLine(line);
chart.render();

`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *AppProcessesSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>命令行匹配关键词</td>
	<td>{{source.cmdlineKeyword}}</td>
</tr>
`
	return p
}

// 平台限制
func (this *AppProcessesSource) Platforms() []string {
	return []string{"darwin", "linux"}
}
