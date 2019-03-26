package teaconfigs

import "github.com/iwind/TeaGo/maps"

// 运算符定义
type RewriteOperator = string

const (
	RewriteOperatorRegexp      = "regexp"
	RewriteOperatorNotRegexp   = "not regexp"
	RewriteOperatorGt          = "gt"
	RewriteOperatorGte         = "gte"
	RewriteOperatorLt          = "lt"
	RewriteOperatorLte         = "lte"
	RewriteOperatorEq          = "eq"
	RewriteOperatorNot         = "not"
	RewriteOperatorPrefix      = "prefix"
	RewriteOperatorSuffix      = "suffix"
	RewriteOperatorContains    = "contains"
	RewriteOperatorNotContains = "not contains"
)

// 所有的运算符
func AllRewriteOperators() []maps.Map {
	return []maps.Map{
		{
			"name":        "正则表达式匹配",
			"op":          RewriteOperatorRegexp,
			"description": "判断是否正则表达式匹配",
		},
		{
			"name":        "正则表达式不匹配",
			"op":          RewriteOperatorNotRegexp,
			"description": "判断是否正则表达式不匹配",
		},
		{
			"name":        "等于",
			"op":          RewriteOperatorEq,
			"description": "使用字符串对比参数值是否相等于某个值",
		},
		{
			"name":        "前缀",
			"op":          RewriteOperatorPrefix,
			"description": "参数值包含某个前缀",
		},
		{
			"name":        "后缀",
			"op":          RewriteOperatorSuffix,
			"description": "参数值包含某个后缀",
		},
		{
			"name":        "包含",
			"op":          RewriteOperatorContains,
			"description": "参数值包含另外一个字符串",
		},
		{
			"name":        "不包含",
			"op":          RewriteOperatorNotContains,
			"description": "参数值不包含另外一个字符串",
		},
		{
			"name":        "不等于",
			"op":          RewriteOperatorNot,
			"description": "使用字符串对比参数值是否不相等于某个值",
		},
		{
			"name":        "大于",
			"op":          RewriteOperatorGt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "大于等于",
			"op":          RewriteOperatorGte,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于",
			"op":          RewriteOperatorLt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于等于",
			"op":          RewriteOperatorLte,
			"description": "将参数转换为数字进行对比",
		},
	}
}

// 查找某个运算符信息
func FindRewriteOperator(op string) maps.Map {
	for _, o := range AllRewriteOperators() {
		if o["op"] == op {
			return o
		}
	}
	return nil
}
