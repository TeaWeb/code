package teaconfigs

// 头部信息定义
// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
type HeaderConfig struct {
	Name   string   `yaml:"name" json:"name"`     // @TODO
	Value  string   `yaml:"value" json:"value"`   // @TODO
	Always bool     `yaml:"always" json:"always"` // @TODO
	Code   []string `yaml:"code" json:"code"`     // @TODO
}
