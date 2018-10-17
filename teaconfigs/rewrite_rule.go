package teaconfigs

import (
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

const (
	RewriteTargetProxy = 1
	RewriteTargetURL   = 2
)

// 重写规则定义
//
// 参考
// - http://nginx.org/en/docs/http/ngx_http_rewrite_module.html
// - https://httpd.apache.org/docs/current/mod/mod_rewrite.html
// - https://httpd.apache.org/docs/2.4/rewrite/flags.html
type RewriteRule struct {
	On bool   `yaml:"on" json:"on"` // 是否开启
	Id string `yaml:"id" json:"id"` // ID

	// 开启的条件
	// 语法为：cond testString condPattern 比如：
	// - cond ${status} 200
	// - cond ${arg.name} lily
	// - cond ${requestPath} *.png
	// @TODO 需要实现
	Cond []RewriteCond `yaml:"cond" json:"cond"`

	// 规则
	// 语法为：pattern regexp 比如：
	// - pattern ^/article/(\d+).html
	Pattern string `yaml:"pattern" json:"pattern"`
	reg     *regexp.Regexp

	// 要替换成的URL
	// 支持反向引用：${0}, ${1}, ...
	// - 如果以 proxy:// 开头，表示目标为代理，首先会尝试作为代理ID请求，如果找不到，会尝试作为代理Host请求
	Replace string `yaml:"replace" json:"replace"`

	// Headers
	Headers       []*HeaderConfig `yaml:"headers" json:"headers"`             // 自定义Header @TODO
	IgnoreHeaders []string        `yaml:"ignoreHeaders" json:"ignoreHeaders"` // 忽略的Header @TODO

	targetType  int // RewriteTarget*
	targetURL   string
	targetProxy string
}

func NewRewriteRule() *RewriteRule {
	return &RewriteRule{
		On: true,
		Id: stringutil.Rand(16),
	}
}

func (this *RewriteRule) Validate() error {
	reg, err := regexp.Compile(this.Pattern)
	if err != nil {
		return err
	}
	this.reg = reg

	// 替换replace中的反向引用
	if strings.HasPrefix(this.Replace, "proxy://") {
		this.targetType = RewriteTargetProxy
		url := this.Replace[len("proxy://"):]
		index := strings.Index(url, "/")
		if index >= 0 {
			this.targetProxy = url[:index]
			this.targetURL = url[index:]
		}
	} else {
		this.targetType = RewriteTargetURL
		this.targetURL = this.Replace
	}

	// 校验条件
	for _, cond := range this.Cond {
		err := cond.Validate()
		if err != nil {
			return err
		}
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

// 对某个请求执行规则
func (this *RewriteRule) Apply(requestPath string, formatter func(source string) string) bool {
	if this.reg == nil {
		return false
	}

	// 判断条件
	for _, cond := range this.Cond {
		if !cond.Match(formatter) {
			return false
		}
	}

	replace := formatter(this.targetURL)
	matches := this.reg.FindStringSubmatch(requestPath)
	if len(matches) == 0 {
		return false
	}
	replace = regexp.MustCompile("\\${\\d+}").ReplaceAllStringFunc(replace, func(s string) string {
		index := types.Int(s[2 : len(s)-1])
		if index < len(matches) {
			return matches[index]
		}
		return ""
	})

	this.targetURL = replace

	return true
}

func (this *RewriteRule) TargetType() int {
	return this.targetType
}

func (this *RewriteRule) TargetProxy() string {
	return this.targetProxy
}

func (this *RewriteRule) TargetURL() string {
	return this.targetURL
}

// 设置Header
func (this *RewriteRule) SetHeader(name string, value string) {
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
func (this *RewriteRule) DeleteHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.Headers) {
		this.Headers = lists.Remove(this.Headers, index).([]*HeaderConfig)
	}
}

// 取得指定位置上的Header
func (this *RewriteRule) HeaderAtIndex(index int) *HeaderConfig {
	if index >= 0 && index < len(this.Headers) {
		return this.Headers[index]
	}
	return nil
}

// 屏蔽一个Header
func (this *RewriteRule) AddIgnoreHeader(name string) {
	this.IgnoreHeaders = append(this.IgnoreHeaders, name)
}

// 移除对Header的屏蔽
func (this *RewriteRule) DeleteIgnoreHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders = lists.Remove(this.IgnoreHeaders, index).([]string)
	}
}

// 更改Header的屏蔽
func (this *RewriteRule) UpdateIgnoreHeaderAtIndex(index int, name string) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders[index] = name
	}
}
