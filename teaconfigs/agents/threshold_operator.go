package agents

import (
	"github.com/iwind/TeaGo/maps"
)

// 运算符定义
type ThresholdOperator = string

const (
	ThresholdOperatorRegexp      = "regexp"
	ThresholdOperatorNotRegexp   = "not regexp"
	ThresholdOperatorGt          = "gt"
	ThresholdOperatorGte         = "gte"
	ThresholdOperatorLt          = "lt"
	ThresholdOperatorLte         = "lte"
	ThresholdOperatorEq          = "eq"
	ThresholdOperatorNumberEq    = "number eq"
	ThresholdOperatorNot         = "not"
	ThresholdOperatorPrefix      = "prefix"
	ThresholdOperatorSuffix      = "suffix"
	ThresholdOperatorContains    = "contains"
	ThresholdOperatorNotContains = "not contains"
)

// 所有的运算符
func AllThresholdOperators() []maps.Map {
	return []maps.Map{
		{
			"name":        "匹配正则表达式",
			"op":          ThresholdOperatorRegexp,
			"description": "判断参数值是否匹配正则表达式",
		},
		{
			"name":        "不匹配正则表达式",
			"op":          ThresholdOperatorNotRegexp,
			"description": "判断参数值是否不匹配正则表达式",
		},
		{
			"name":        "字符串等于",
			"op":          ThresholdOperatorEq,
			"description": "使用字符串对比参数值是否相等于某个值",
		},
		{
			"name":        "数字等于",
			"op":          ThresholdOperatorNumberEq,
			"description": "使用数字对比参数值是否相等于某个值",
		},
		{
			"name":        "不等于",
			"op":          ThresholdOperatorNot,
			"description": "使用字符串对比参数值是否不相等于某个值",
		},
		{
			"name":        "前缀",
			"op":          ThresholdOperatorPrefix,
			"description": "参数值包含某个前缀",
		},
		{
			"name":        "后缀",
			"op":          ThresholdOperatorSuffix,
			"description": "参数值包含某个后缀",
		},
		{
			"name":        "包含",
			"op":          ThresholdOperatorContains,
			"description": "参数值包含另外一个字符串",
		},
		{
			"name":        "不包含",
			"op":          ThresholdOperatorNotContains,
			"description": "参数值不包含另外一个字符串",
		},
		{
			"name":        "大于",
			"op":          ThresholdOperatorGt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "大于等于",
			"op":          ThresholdOperatorGte,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于",
			"op":          ThresholdOperatorLt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于等于",
			"op":          ThresholdOperatorLte,
			"description": "将参数转换为数字进行对比",
		},
	}
}

// 查找某个运算符信息
func FindThresholdOperator(op string) maps.Map {
	for _, o := range AllThresholdOperators() {
		if o["op"] == op {
			return o
		}
	}
	return nil
}
