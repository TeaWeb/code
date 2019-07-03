package agents

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"time"
)

// URL连通性
type URLConnectivitySource struct {
	Source `yaml:",inline"`

	Timeout int    `yaml:"timeout" json:"timeout"` // 连接超时
	URL     string `yaml:"url" json:"url"`
	Method  string `yaml:"method" json:"method"`
}

// 获取新对象
func NewURLConnectivitySource() *URLConnectivitySource {
	return &URLConnectivitySource{}
}

// 名称
func (this *URLConnectivitySource) Name() string {
	return "URL连通性"
}

// 代号
func (this *URLConnectivitySource) Code() string {
	return "urlConnectivity"
}

// 描述
func (this *URLConnectivitySource) Description() string {
	return "获取URL连通性信息"
}

// 执行
func (this *URLConnectivitySource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.URL) == 0 {
		err = errors.New("'url' should not be empty")
		return maps.Map{
			"status": 0,
		}, err
	}

	method := this.Method
	if len(method) == 0 {
		method = http.MethodGet
	}

	before := time.Now()
	req, err := http.NewRequest(method, this.URL, nil)
	if err != nil {
		value = maps.Map{
			"cost":   time.Since(before).Seconds(),
			"status": 0,
			"result": "",
			"length": 0,
		}
		return value, err
	}
	req.Header.Set("User-Agent", "TeaWeb/"+teaconst.TeaVersion)

	timeout := this.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := teautils.NewHttpClient(time.Duration(timeout) * time.Second)
	defer teautils.CloseHTTPClient(client)

	resp, err := client.Do(req)
	if err != nil {
		return maps.Map{
			"status": 0,
		}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return maps.Map{
			"status": 0,
		}, err
	}

	if len(data) > 1024 {
		data = data[:1024]
	}

	value = maps.Map{
		"cost":   time.Since(before).Seconds(),
		"status": resp.StatusCode,
		"result": string(data),
		"length": len(data),
	}

	return
}

// 表单信息
func (this *URLConnectivitySource) Form() *forms.Form {
	form := forms.NewForm(this.Code())

	{
		group := form.NewGroup()

		{
			field := forms.NewTextField("URL", "")
			field.Code = "url"
			field.IsRequired = true
			field.Placeholder = "http://..."
			field.MaxLength = 500
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入URL")
}

if (!value.match(/^(http|https):\/\//i)) {
	throw new Error("URL地址必须以http或https开头");
}
`
			group.Add(field)
		}

		{
			field := forms.NewOptions("请求方法", "")
			field.Code = "method"
			field.IsRequired = true
			field.AddOption("GET", "GET")
			field.AddOption("POST", "POST")
			field.AddOption("PUT", "PUT")
			field.Attr("style", "width:10em")
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请选择请求方法");
}
`
			group.Add(field)
		}
	}

	{
		group := form.NewGroup()

		{
			field := forms.NewTextField("请求超时", "Timeout")
			field.Code = "timeout"
			field.Value = 10
			field.MaxLength = 10
			field.RightLabel = "秒"
			field.Attr("style", "width:5em")
			field.ValidateCode = `
var intValue = parseInt(value);
if (isNaN(intValue)) {
	throw new Error("超时时间需要是一个整数");
}

return intValue;
`

			group.Add(field)
		}
	}

	return form
}

// 变量
func (this *URLConnectivitySource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "cost",
			Description: "请求耗时（秒）",
		},
		{
			Code:        "status",
			Description: "HTTP状态码",
		},
		{
			Code:        "result",
			Description: "响应内容文本，最多只记录前1024个字节",
		},
		{
			Code:        "length",
			Description: "响应的内容长度",
		},
	}
}

// 阈值
func (this *URLConnectivitySource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${status}"
		t.Operator = ThresholdOperatorGte
		t.Value = "400"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "请求状态码错误"
		result = append(result, t)
	}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${status}"
		t.Operator = ThresholdOperatorEq
		t.Value = "0"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "URL请求失败"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *URLConnectivitySource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "URL连通性（ms）"
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
	if (v.value == "") {
		return;
	}
	line.values.push(v.value.cost * 1000);
	
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
func (this *URLConnectivitySource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>URL</td>
	<td>{{source.url}}</td>
</tr>
<tr>
	<td>请求方法</td>
	<td>{{source.method}}</td>
</tr>
<tr>
	<td>请求超时</td>
	<td>{{source.timeout}}s</td>
</tr>
`
	return p
}
