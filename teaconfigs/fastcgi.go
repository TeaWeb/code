package teaconfigs

import (
	"github.com/iwind/TeaGo/maps"
	"github.com/TeaWeb/code/teaconst"
	"net/http"
	"path/filepath"
	"github.com/iwind/TeaGo/utils/string"
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
