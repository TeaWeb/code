package teaconfigs

import (
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

// 重写条件定义
type RewriteCond struct {
	Id string `yaml:"id" json:"id"` // ID

	// 要测试的字符串
	// 其中可以使用跟请求相关的参数，比如：
	// ${arg.name}, ${requestPath}
	Param string `yaml:"param" json:"param"`

	// 运算符
	Operator RewriteOperator `yaml:"operator" json:"operator"`

	// 对比
	Value string `yaml:"value" json:"value"`

	regValue   *regexp.Regexp
	floatValue float64
}

// 取得新对象
func NewRewriteCond() *RewriteCond {
	return &RewriteCond{
		Id: stringutil.Rand(16),
	}
}

// 校验配置
func (this *RewriteCond) Validate() error {
	if this.Operator == RewriteOperatorRegexp {
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.regValue = reg
	} else if this.Operator == RewriteOperatorGt || this.Operator == RewriteOperatorGte || this.Operator == RewriteOperatorLt || this.Operator == RewriteOperatorLte {
		this.floatValue = types.Float64(this.Value)
	}
	return nil
}

// 将此条件应用于请求，检查是否匹配
func (this *RewriteCond) Match(formatter func(source string) string) bool {
	paramValue := formatter(this.Param)
	switch this.Operator {
	case RewriteOperatorRegexp:
		if this.regValue == nil {
			return false
		}
		return this.regValue.MatchString(paramValue)
	case RewriteOperatorGt:
		return types.Float64(paramValue) > this.floatValue
	case RewriteOperatorGte:
		return types.Float64(paramValue) >= this.floatValue
	case RewriteOperatorLt:
		return types.Float64(paramValue) < this.floatValue
	case RewriteOperatorLte:
		return types.Float64(paramValue) <= this.floatValue
	case RewriteOperatorEq:
		return paramValue == this.Value
	case RewriteOperatorNot:
		return paramValue != this.Value
	case RewriteOperatorPrefix:
		return strings.HasPrefix(paramValue, this.Value)
	case RewriteOperatorSuffix:
		return strings.HasSuffix(paramValue, this.Value)
	case RewriteOperatorContains:
		return strings.Contains(paramValue, this.Value)
	}
	return false
}
