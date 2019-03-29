package teaconfigs

import "github.com/iwind/TeaGo/maps"

// 运算符定义
type RequestCondOperator = string

const (
	RequestCondOperatorRegexp      = "regexp"
	RequestCondOperatorNotRegexp   = "not regexp"
	RequestCondOperatorGt          = "gt"
	RequestCondOperatorGte         = "gte"
	RequestCondOperatorLt          = "lt"
	RequestCondOperatorLte         = "lte"
	RequestCondOperatorEq          = "eq"
	RequestCondOperatorNot         = "not"
	RequestCondOperatorPrefix      = "prefix"
	RequestCondOperatorSuffix      = "suffix"
	RequestCondOperatorContains    = "contains"
	RequestCondOperatorNotContains = "not contains"
)

// 所有的运算符
func AllRequestOperators() []maps.Map {
	return []maps.Map{
		{
			"name":        "正则表达式匹配",
			"op":          RequestCondOperatorRegexp,
			"description": "判断是否正则表达式匹配",
		},
		{
			"name":        "正则表达式不匹配",
			"op":          RequestCondOperatorNotRegexp,
			"description": "判断是否正则表达式不匹配",
		},
		{
			"name":        "等于",
			"op":          RequestCondOperatorEq,
			"description": "使用字符串对比参数值是否相等于某个值",
		},
		{
			"name":        "前缀",
			"op":          RequestCondOperatorPrefix,
			"description": "参数值包含某个前缀",
		},
		{
			"name":        "后缀",
			"op":          RequestCondOperatorSuffix,
			"description": "参数值包含某个后缀",
		},
		{
			"name":        "包含",
			"op":          RequestCondOperatorContains,
			"description": "参数值包含另外一个字符串",
		},
		{
			"name":        "不包含",
			"op":          RequestCondOperatorNotContains,
			"description": "参数值不包含另外一个字符串",
		},
		{
			"name":        "不等于",
			"op":          RequestCondOperatorNot,
			"description": "使用字符串对比参数值是否不相等于某个值",
		},
		{
			"name":        "大于",
			"op":          RequestCondOperatorGt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "大于等于",
			"op":          RequestCondOperatorGte,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于",
			"op":          RequestCondOperatorLt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于等于",
			"op":          RequestCondOperatorLte,
			"description": "将参数转换为数字进行对比",
		},
	}
}

// 查找某个运算符信息
func FindRequestCondOperator(op string) maps.Map {
	for _, o := range AllRequestOperators() {
		if o["op"] == op {
			return o
		}
	}
	return nil
}
