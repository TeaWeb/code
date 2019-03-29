package teaconfigs

import (
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

// 重写条件定义
type RequestCond struct {
	Id string `yaml:"id" json:"id"` // ID

	// 要测试的字符串
	// 其中可以使用跟请求相关的参数，比如：
	// ${arg.name}, ${requestPath}
	Param string `yaml:"param" json:"param"`

	// 运算符
	Operator RequestCondOperator `yaml:"operator" json:"operator"`

	// 对比
	Value string `yaml:"value" json:"value"`

	regValue   *regexp.Regexp
	floatValue float64
}

// 取得新对象
func NewRequestCond() *RequestCond {
	return &RequestCond{
		Id: stringutil.Rand(16),
	}
}

// 校验配置
func (this *RequestCond) Validate() error {
	if this.Operator == RequestCondOperatorRegexp || this.Operator == RequestCondOperatorNotRegexp {
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.regValue = reg
	} else if this.Operator == RequestCondOperatorGt || this.Operator == RequestCondOperatorGte || this.Operator == RequestCondOperatorLt || this.Operator == RequestCondOperatorLte {
		this.floatValue = types.Float64(this.Value)
	}
	return nil
}

// 将此条件应用于请求，检查是否匹配
func (this *RequestCond) Match(formatter func(source string) string) bool {
	paramValue := formatter(this.Param)
	switch this.Operator {
	case RequestCondOperatorRegexp:
		if this.regValue == nil {
			return false
		}
		return this.regValue.MatchString(paramValue)
	case RequestCondOperatorNotRegexp:
		if this.regValue == nil {
			return false
		}
		return !this.regValue.MatchString(paramValue)
	case RequestCondOperatorGt:
		return types.Float64(paramValue) > this.floatValue
	case RequestCondOperatorGte:
		return types.Float64(paramValue) >= this.floatValue
	case RequestCondOperatorLt:
		return types.Float64(paramValue) < this.floatValue
	case RequestCondOperatorLte:
		return types.Float64(paramValue) <= this.floatValue
	case RequestCondOperatorEq:
		return paramValue == this.Value
	case RequestCondOperatorNot:
		return paramValue != this.Value
	case RequestCondOperatorPrefix:
		return strings.HasPrefix(paramValue, this.Value)
	case RequestCondOperatorSuffix:
		return strings.HasSuffix(paramValue, this.Value)
	case RequestCondOperatorContains:
		return strings.Contains(paramValue, this.Value)
	case RequestCondOperatorNotContains:
		return !strings.Contains(paramValue, this.Value)
	}
	return false
}
