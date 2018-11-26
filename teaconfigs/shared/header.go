package shared

import "github.com/iwind/TeaGo/lists"

// 头部信息定义
// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
type HeaderConfig struct {
	Name   string `yaml:"name" json:"name"`     // Name
	Value  string `yaml:"value" json:"value"`   // Value
	Always bool   `yaml:"always" json:"always"` // 是否忽略状态码 @TODO
	Status []int  `yaml:"code" json:"code"`     // 支持的状态码 @TODO
}

func NewHeaderConfig() *HeaderConfig {
	return &HeaderConfig{}
}

func (this *HeaderConfig) Validate() error {
	return nil
}

func (this *HeaderConfig) Match(statusCode int) bool {
	if this.Always {
		return true
	}

	if len(this.Status) > 0 {
		return lists.Contains(this.Status, statusCode)
	}

	return lists.Contains([]int{200, 201, 204, 206, 301, 302, 303, 304, 307, 308}, statusCode)
}
