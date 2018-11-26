package teaconfigs

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Fastcgi配置
// 参考：http://nginx.org/en/docs/http/ngx_http_fastcgi_module.html
type FastcgiConfig struct {
	On bool   `yaml:"on" json:"on"` // @TODO
	Id string `yaml:"id" json:"id"` // @TODO

	// fastcgi地址配置
	// 支持unix:/tmp/php-fpm.sock ...
	Pass string `yaml:"pass" json:"pass"`

	Index       string            `yaml:"index" json:"index"`             //@TODO
	Params      map[string]string `yaml:"params" json:"params"`           //@TODO
	ReadTimeout string            `yaml:"readTimeout" json:"readTimeout"` // @TODO 读取超时时间
	SendTimeout string            `yaml:"sendTimeout" json:"sendTimeout"` // @TODO 发送超时时间
	ConnTimeout string            `yaml:"connTimeout" json:"connTimeout"` // @TODO 连接超时时间
	PoolSize    int               `yaml:"poolSize" json:"poolSize"`       // 连接池尺寸 @TODO

	// Headers
	Headers       []*shared.HeaderConfig `yaml:"headers" json:"headers"`             // 自定义Header @TODO
	IgnoreHeaders []string               `yaml:"ignoreHeaders" json:"ignoreHeaders"` // 忽略的Header @TODO

	network string // 协议：tcp, unix
	address string // 地址

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

	// 校验地址
	if regexp.MustCompile("^\\d+$").MatchString(this.Pass) {
		this.network = "tcp"
		this.address = "127.0.0.1:" + this.Pass
	} else if regexp.MustCompile("^(.*):(\\d+)$").MatchString(this.Pass) {
		matches := regexp.MustCompile("^(.*):(\\d+)$").FindStringSubmatch(this.Pass)
		ip := matches[1]
		port := matches[2]
		if len(ip) == 0 {
			ip = "127.0.0.1"
		}
		this.network = "tcp"
		this.address = ip + ":" + port
	} else if regexp.MustCompile("^\\d+\\.\\d+.\\d+.\\d+$").MatchString(this.Pass) {
		this.network = "tcp"
		this.address = this.Pass + ":9000"
	} else if regexp.MustCompile("^unix:(.+)$").MatchString(this.Pass) {
		matches := regexp.MustCompile("^unix:(.+)$").FindStringSubmatch(this.Pass)
		path := matches[1]
		this.network = "unix"
		this.address = path
	} else if regexp.MustCompile("^[./].+$").MatchString(this.Pass) {
		this.network = "unix"
		this.address = this.Pass
	} else {
		return errors.New("invalid 'pass' format")
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

//网络协议
func (this *FastcgiConfig) Network() string {
	return this.network
}

// 网络地址
func (this *FastcgiConfig) Address() string {
	return this.address
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

	header := shared.NewHeaderConfig()
	header.Name = name
	header.Value = value
	this.Headers = append(this.Headers, header)
}

// 删除指定位置上的Header
func (this *FastcgiConfig) DeleteHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.Headers) {
		this.Headers = lists.Remove(this.Headers, index).([]*shared.HeaderConfig)
	}
}

// 取得指定位置上的Header
func (this *FastcgiConfig) HeaderAtIndex(index int) *shared.HeaderConfig {
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
