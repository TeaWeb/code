package agents

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
	"github.com/tatsushid/go-fastping"
	"net"
	"runtime"
	"time"
)

// Ping
type PingSource struct {
	Source `yaml:",inline"`

	Host string `yaml:"host" json:"host"`
}

// 获取新对象
func NewPingSource() *PingSource {
	return &PingSource{}
}

// 名称
func (this *PingSource) Name() string {
	return "Ping"
}

// 代号
func (this *PingSource) Code() string {
	return "ping"
}

// 描述
func (this *PingSource) Description() string {
	return "通过Ping获取主机响应时间"
}

// 执行
func (this *PingSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.Host) == 0 {
		err = errors.New("'host' should not be empty")
		return maps.Map{
			"rtt": -1,
		}, err
	}

	p := fastping.NewPinger()
	if runtime.GOOS == "darwin" {
		_, err = p.Network("udp")
	} else {
		_, err = p.Network("ip")
	}
	if err != nil {
		return maps.Map{
			"rtt": -1,
		}, err
	}

	ra, err := net.ResolveIPAddr("ip4:icmp", this.Host)
	if err != nil {
		return maps.Map{
			"rtt": -1,
		}, err
	}
	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		value = maps.Map{
			"rtt": rtt.Seconds() * 1000,
		}
	}
	p.OnIdle = func() {
		if value == nil {
			err = errors.New("ping timeout")
		}
	}

	runningErr := p.Run()
	if runningErr != nil {
		return maps.Map{
			"rtt": -1,
		}, runningErr
	}

	if err != nil {
		return maps.Map{
			"rtt": -1,
		}, err
	}

	return
}

// 表单信息
func (this *PingSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewTextField("主机地址", "Host")
			field.IsRequired = true
			field.Code = "host"
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入主机地址");
}
`
			field.Comment = "要Ping的主机地址，可以是一个域名或一个IP"
			group.Add(field)
		}
	}
	return form
}

// 变量
func (this *PingSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "rtt",
			Description: "响应时间（单位ms）",
		},
	}
}

// 阈值
func (this *PingSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${rtt}"
		t.Operator = ThresholdOperatorEq
		t.Value = "-1"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "Ping超时"
		t.MaxFails = 5
		result = append(result, t)
	}

	return result
}

// 图表
func (this *PingSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "Ping"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("rtt");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.rtt);
	chart.addLabel(v.label);
});

chart.addLine(line);
chart.render();`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *PingSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>主机地址</td>
	<td>{{source.host}}</td>
</tr>
`
	return p
}
