package teaconfigs

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Fastcgi配置
// 参考：http://nginx.org/en/docs/http/ngx_http_fastcgi_module.html
type FastcgiConfig struct {
	On bool   `yaml:"on" json:"on"` // @TODO
	Id string `yaml:"id" json:"id"` // @TODO

	// @TODO 支持unix://...
	Pass string `yaml:"pass" json:"pass"` // @TODO

	Index       string            `yaml:"index" json:"index"`             //@TODO
	Params      map[string]string `yaml:"params" json:"params"`           //@TODO
	ReadTimeout string            `yaml:"readTimeout" json:"readTimeout"` // @TODO 读取超时时间
	SendTimeout string            `yaml:"sendTimeout" json:"sendTimeout"` // @TODO 发送超时时间
	ConnTimeout string            `yaml:"connTimeout" json:"connTimeout"` // @TODO 连接超时时间
	PoolSize    int               `yaml:"poolSize" json:"poolSize"`       // 连接池尺寸 @TODO

	// Headers
	Headers       []*HeaderConfig `yaml:"headers" json:"headers"`             // 自定义Header @TODO
	IgnoreHeaders []string        `yaml:"ignoreHeaders" json:"ignoreHeaders"` // 忽略的Header @TODO

	paramsMap maps.Map
	timeout   time.Duration
}

func NewFastcgiConfig() *FastcgiConfig {
	return &FastcgiConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验配置
func (this *FastcgiConfig) Validate() error {
	this.paramsMap = maps.NewMap(this.Params)
	if !this.paramsMap.Has("SCRIPT_FILENAME") {
		this.paramsMap["SCRIPT_FILENAME"] = ""
	}
	if !this.paramsMap.Has("SERVER_SOFTWARE") {
		this.paramsMap["SERVER_SOFTWARE"] = "teaweb/" + teaconst.TeaVersion
	}
	if !this.paramsMap.Has("REDIRECT_STATUS") {
		this.paramsMap["REDIRECT_STATUS"] = "200"
	}
	if !this.paramsMap.Has("GATEWAY_INTERFACE") {
		this.paramsMap["GATEWAY_INTERFACE"] = "CGI/1.1"
	}

	// 超时时间
	if len(this.ReadTimeout) > 0 {
		duration, err := time.ParseDuration(this.ReadTimeout)
		if err != nil {
			return err
		}
		this.timeout = duration
	} else {
		this.timeout = 3 * time.Second
	}

	// 校验Header
	for _, header := range this.Headers {
		err := header.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *FastcgiConfig) FilterParams(req *http.Request) maps.Map {
	params := maps.NewMap(this.paramsMap)

	//@TODO 处理参数中的${varName}变量

	// 自动添加参数
	script := params.GetString("SCRIPT_FILENAME")
	if len(script) > 0 {
		if !params.Has("SCRIPT_NAME") {
			params["SCRIPT_NAME"] = filepath.Base(script)
		}
		if !params.Has("DOCUMENT_ROOT") {
			params["DOCUMENT_ROOT"] = filepath.Dir(script)
		}
		if !params.Has("PWD") {
			params["PWD"] = filepath.Dir(script)
		}
	}

	return params
}

// 超时时间
func (this *FastcgiConfig) Timeout() time.Duration {
	if this.timeout <= 0 {
		this.timeout = 30 * time.Second
	}
	return this.timeout
}

// 设置Header
func (this *FastcgiConfig) SetHeader(name string, value string) {
	found := false
	upperName := strings.ToUpper(name)
	for _, header := range this.Headers {
		if strings.ToUpper(header.Name) == upperName {
			found = true
			header.Value = value
		}
	}
	if found {
		return
	}

	header := NewHeaderConfig()
	header.Name = name
	header.Value = value
	this.Headers = append(this.Headers, header)
}

// 删除指定位置上的Header
func (this *FastcgiConfig) DeleteHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.Headers) {
		this.Headers = lists.Remove(this.Headers, index).([]*HeaderConfig)
	}
}

// 取得指定位置上的Header
func (this *FastcgiConfig) HeaderAtIndex(index int) *HeaderConfig {
	if index >= 0 && index < len(this.Headers) {
		return this.Headers[index]
	}
	return nil
}

// 屏蔽一个Header
func (this *FastcgiConfig) AddIgnoreHeader(name string) {
	this.IgnoreHeaders = append(this.IgnoreHeaders, name)
}

// 移除对Header的屏蔽
func (this *FastcgiConfig) DeleteIgnoreHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders = lists.Remove(this.IgnoreHeaders, index).([]string)
	}
}

// 更改Header的屏蔽
func (this *FastcgiConfig) UpdateIgnoreHeaderAtIndex(index int, name string) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders[index] = name
	}
}
