package agents

import (
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/net"
	"time"
)

// 网络带宽等信息
type NetworkSource struct {
	Source `yaml:",inline"`

	lastTimer time.Time

	countSent     uint64
	countReceived uint64

	lastTotalSent     uint64
	lastTotalReceived uint64
}

// 获取新对象
func NewNetworkSource() *NetworkSource {
	return &NetworkSource{}
}

// 名称
func (this *NetworkSource) Name() string {
	return "网络信息"
}

// 代号
func (this *NetworkSource) Code() string {
	return "network"
}

// 描述
func (this *NetworkSource) Description() string {
	return "网络接口、带宽信息"
}

// 执行
func (this *NetworkSource) Execute(params map[string]string) (value interface{}, err error) {
	interfaces := []map[string]interface{}{}
	interfaceStats, err := net.Interfaces()
	if err != nil {
		logs.Error(err)
	} else {
		for _, i := range interfaceStats {
			interfaces = append(interfaces, map[string]interface{}{
				"name":         i.Name,
				"hardwareAddr": i.HardwareAddr,
				"mtu":          i.MTU,
				"flags":        i.Flags,
				"addrs": lists.Map(i.Addrs, func(k int, v interface{}) interface{} {
					addr := v.(net.InterfaceAddr)
					return addr.Addr
				}),
			})
		}
	}

	value = map[string]interface{}{
		"interfaces": interfaces,
		"stat": map[string]interface{}{
			"avgSentBytes":       0,
			"avgReceivedBytes":   0,
			"totalSentBytes":     0,
			"totalReceivedBytes": 0,
		},
	}

	stats, err := net.IOCounters(false)
	if err != nil {
		return value, err
	}
	totalSent := uint64(0)
	totalReceived := uint64(0)
	for _, stat := range stats {
		totalSent += stat.BytesSent
		totalReceived += stat.BytesRecv
	}

	if this.lastTotalSent == 0 {
		this.lastTotalSent = totalSent
		this.lastTotalReceived = totalReceived
		return
	}

	if totalSent > this.lastTotalSent {
		this.countSent = totalSent - this.lastTotalSent
	} else {
		this.countSent = 0
	}
	if totalReceived > this.lastTotalReceived {
		this.countReceived = totalReceived - this.lastTotalReceived
	} else {
		this.countReceived = 0
	}

	this.lastTotalSent = totalSent
	this.lastTotalReceived = totalReceived

	duration := time.Since(this.lastTimer).Seconds()
	if duration <= 0 {
		return
	}

	value = map[string]interface{}{
		"interfaces": interfaces,
		"stat": map[string]interface{}{
			"avgSentBytes":       this.countSent / uint64(duration),
			"avgReceivedBytes":   this.countReceived / uint64(duration),
			"totalSentBytes":     totalSent,
			"totalReceivedBytes": totalReceived,
		},
	}

	this.lastTimer = time.Now()

	return
}

// 表单信息
func (this *NetworkSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *NetworkSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "interfaces",
			Description: "网络接口",
		},
		{
			Code:        "interfaces.$.name",
			Description: "接口名称",
		},
		{
			Code:        "interfaces.$.addrs",
			Description: "接口地址",
		},
		{
			Code:        "interfaces.$.flags",
			Description: "接口标识",
		},
		{
			Code:        "interfaces.$.hardwareAddr",
			Description: "接口硬件地址",
		},
		{
			Code:        "interfaces.$.mtu",
			Description: "接口MTU值",
		},
		{
			Code:        "stat",
			Description: "流量统计信息",
		},
		{
			Code:        "stat.avgReceivedBytes",
			Description: "平均接收速率（秒）",
		},
		{
			Code:        "stat.avgSentBytes",
			Description: "平均发送速率（秒）",
		},
		{
			Code:        "stat.totalReceivedBytes",
			Description: "总接收字节数",
		},
		{
			Code:        "stat.totalSentBytes",
			Description: "总发送字节数",
		},
	}
}

// 阈值
func (this *NetworkSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${stat.avgSentBytes}"
		t.Operator = ThresholdOperatorGte
		t.Value = "13107200"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "当前出口流量超过100MBit/s"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *NetworkSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	// 图表
	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.received"
		chart.Name = "出口带宽（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = new values.Query().cache(60).latest(60);
ones.reverse();
ones.$each(function (k, v) {
	line.values.push(Math.round(v.value.stat.avgSentBytes / 1024 / 1024 * 100) / 100);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();
`,
		}
		charts = append(charts, chart)
	}

	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.sent"
		chart.Name = "入口带宽（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = new values.Query().cache(60).latest(60);
ones.reverse();
ones.$each(function (k, v) {
	line.values.push(Math.round(v.value.stat.avgReceivedBytes / 1024 / 1024 * 100) / 100);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();
`,
		}
		charts = append(charts, chart)
	}

	return charts
}
