package agents

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teautils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebHook
type WebHookSource struct {
	Source `yaml:",inline"`

	URL     string `yaml:"url" json:"url"`
	Timeout string `yaml:"timeout" json:"timeout"`
	Method  string `yaml:"method" json:"method"` // 请求方法

	timeoutDuration time.Duration
}

// 获取新对象
func NewWebHookSource() *WebHookSource {
	return &WebHookSource{}
}

// 校验
func (this *WebHookSource) Validate() error {
	this.timeoutDuration, _ = time.ParseDuration(this.Timeout)
	if len(this.Method) == 0 {
		this.Method = http.MethodPost
	} else {
		this.Method = strings.ToUpper(this.Method)
	}

	if len(this.URL) == 0 {
		return errors.New("url should not be empty")
	}

	return nil
}

// 名称
func (this *WebHookSource) Name() string {
	return "WebHook"
}

// 代号
func (this *WebHookSource) Code() string {
	return "webhook"
}

// 描述
func (this *WebHookSource) Description() string {
	return "通过HTTP或者HTTPS接口获取数据"
}

// 执行
func (this *WebHookSource) Execute(params map[string]string) (value interface{}, err error) {
	if this.timeoutDuration.Seconds() <= 0 {
		this.timeoutDuration = 10 * time.Second
	}

	client := teautils.NewHttpClient(this.timeoutDuration)
	defer teautils.CloseHTTPClient(client)

	query := url.Values{}
	for name, value := range params {
		query[name] = []string{value}
	}
	rawQuery := query.Encode()

	urlString := this.URL
	var body io.Reader = nil
	if len(rawQuery) > 0 {
		if this.Method == http.MethodGet {
			if strings.Index(this.URL, "?") > 0 {
				urlString += "&" + rawQuery
			} else {
				urlString += "?" + rawQuery
			}
		} else {
			body = bytes.NewReader([]byte(rawQuery))
		}
	}

	req, err := http.NewRequest(this.Method, urlString, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "TeaWeb/"+teaconst.TeaVersion)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code should be 200, now is " + fmt.Sprintf("%d", resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return DecodeSource(respBytes, this.DataFormat)
}

// 选项表单
func (this *WebHookSource) Form() *forms.Form {
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

return intValue + "s"
`
			field.InitCode = `
return value.replace("s", "");
`
			group.Add(field)
		}
	}

	return form
}

func (this *WebHookSource) Presentation() *forms.Presentation {
	return &forms.Presentation{
		HTML: `
<tr>
	<td>URL</td>
	<td>{{source.url}}</td>
</tr>
<tr>
	<td>请求方法</td>
	<td>{{source.method}}</td>
</tr>
<tr>
	<td>请求超时<em>（Timeout）</em></td>
	<td>{{source.timeout}}</td>
</tr>`,
	}
}
