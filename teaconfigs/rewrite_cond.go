package teaconfigs

import (
	"regexp"
)

// 重写条件定义
type RewriteCond struct {
	// 要测试的字符串
	// 其中可以使用跟请求相关的参数，比如：
	// ${arg.name}, ${requestPath}
	Test string `yaml:"test" json:"test"`

	// 规则
	Pattern string `yaml:"pattern" json:"pattern"`
	reg     *regexp.Regexp
}

// 校验配置
func (this *RewriteCond) Validate() error {
	reg, err := regexp.Compile(this.Pattern)
	if err != nil {
		return err
	}
	this.reg = reg
	return nil
}

// 将此条件应用于请求，检查是否匹配
func (this *RewriteCond) Match(formatter func(source string) string) bool {
	if this.reg == nil {
		return false
	}

	return this.reg.MatchString(formatter(this.Test))
}
